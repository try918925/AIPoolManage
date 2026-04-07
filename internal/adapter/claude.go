package adapter

import (
	"awesomeProject/internal/model"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ClaudeAdapter struct{}

func (a *ClaudeAdapter) ProviderType() string { return "claude" }

func (a *ClaudeAdapter) ConvertRequest(req *model.ChatRequest, channel *model.ProviderModel, provider *model.Provider, apiKey string) (*http.Request, error) {
	var systemMsg string
	var messages []map[string]string

	for _, msg := range req.Messages {
		if msg.Role == "system" {
			systemMsg = msg.Content
		} else {
			messages = append(messages, map[string]string{
				"role":    msg.Role,
				"content": msg.Content,
			})
		}
	}

	body := map[string]interface{}{
		"model":    channel.ModelID,
		"messages": messages,
		"stream":   req.Stream,
	}

	if systemMsg != "" {
		body["system"] = systemMsg
	}

	maxTokens := 4096
	if req.MaxTokens != nil {
		maxTokens = *req.MaxTokens
	}
	body["max_tokens"] = maxTokens

	if req.Temperature != nil {
		body["temperature"] = *req.Temperature
	}
	if req.TopP != nil {
		body["top_p"] = *req.TopP
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := strings.TrimRight(provider.BaseURL, "/") + "/v1/messages"
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	return httpReq, nil
}

type claudeResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model     string `json:"model"`
	StopReason string `json:"stop_reason"`
	Usage     struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

func (a *ClaudeAdapter) ConvertResponse(resp *http.Response, requestModel string) (*model.ChatResponse, *model.ChatUsage, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("claude api error (status %d): %s", resp.StatusCode, string(body))
	}

	var cr claudeResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return nil, nil, err
	}

	var content string
	for _, c := range cr.Content {
		if c.Type == "text" {
			content += c.Text
		}
	}

	finishReason := "stop"
	if cr.StopReason == "max_tokens" {
		finishReason = "length"
	}

	chatResp := &model.ChatResponse{
		ID:      generateChatID(),
		Object:  "chat.completion",
		Created: nowUnix(),
		Model:   requestModel,
		Choices: []model.ChatChoice{
			{
				Index: 0,
				Message: &model.ChatMessage{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: &finishReason,
			},
		},
		Usage: &model.ChatUsage{
			PromptTokens:     cr.Usage.InputTokens,
			CompletionTokens: cr.Usage.OutputTokens,
			TotalTokens:      cr.Usage.InputTokens + cr.Usage.OutputTokens,
		},
	}

	return chatResp, chatResp.Usage, nil
}

func (a *ClaudeAdapter) ConvertStreamResponse(resp *http.Response, writer http.ResponseWriter, flusher http.Flusher, requestModel string) (*model.ChatUsage, error) {
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("claude stream error (status %d): %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	chatID := generateChatID()
	created := nowUnix()
	usage := &model.ChatUsage{}
	sentRole := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "event: ") {
			eventType := strings.TrimPrefix(line, "event: ")
			if eventType == "message_stop" {
				// Send finish chunk
				chunk := map[string]interface{}{
					"id":      chatID,
					"object":  "chat.completion.chunk",
					"created": created,
					"model":   requestModel,
					"choices": []map[string]interface{}{
						{
							"index":         0,
							"delta":         map[string]interface{}{},
							"finish_reason": "stop",
						},
					},
				}
				data, _ := json.Marshal(chunk)
				fmt.Fprintf(writer, "data: %s\n\n", string(data))
				flusher.Flush()

				fmt.Fprintf(writer, "data: [DONE]\n\n")
				flusher.Flush()
				break
			}
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		var event map[string]interface{}
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		eventType, _ := event["type"].(string)

		switch eventType {
		case "content_block_delta":
			delta, ok := event["delta"].(map[string]interface{})
			if !ok {
				continue
			}
			text, _ := delta["text"].(string)

			chunkDelta := map[string]interface{}{"content": text}
			if !sentRole {
				chunkDelta["role"] = "assistant"
				sentRole = true
			}

			chunk := map[string]interface{}{
				"id":      chatID,
				"object":  "chat.completion.chunk",
				"created": created,
				"model":   requestModel,
				"choices": []map[string]interface{}{
					{
						"index":         0,
						"delta":         chunkDelta,
						"finish_reason": nil,
					},
				},
			}
			out, _ := json.Marshal(chunk)
			fmt.Fprintf(writer, "data: %s\n\n", string(out))
			flusher.Flush()

		case "message_delta":
			if u, ok := event["usage"].(map[string]interface{}); ok {
				if v, ok := u["output_tokens"].(float64); ok {
					usage.CompletionTokens = int(v)
				}
			}

		case "message_start":
			if msg, ok := event["message"].(map[string]interface{}); ok {
				if u, ok := msg["usage"].(map[string]interface{}); ok {
					if v, ok := u["input_tokens"].(float64); ok {
						usage.PromptTokens = int(v)
					}
				}
			}
		}
	}

	usage.TotalTokens = usage.PromptTokens + usage.CompletionTokens
	return usage, scanner.Err()
}
