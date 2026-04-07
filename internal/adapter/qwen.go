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

type QwenAdapter struct{}

func (a *QwenAdapter) ProviderType() string { return "qwen" }

func (a *QwenAdapter) ConvertRequest(req *model.ChatRequest, channel *model.ProviderModel, provider *model.Provider, apiKey string) (*http.Request, error) {
	// Qwen supports OpenAI-compatible format at /v1/chat/completions
	// Use the compatible endpoint for simplicity
	body := map[string]interface{}{
		"model":    channel.ModelID,
		"messages": req.Messages,
		"stream":   req.Stream,
	}
	if req.Temperature != nil {
		body["temperature"] = *req.Temperature
	}
	if req.MaxTokens != nil {
		body["max_tokens"] = *req.MaxTokens
	}
	if req.TopP != nil {
		body["top_p"] = *req.TopP
	}
	if req.Stop != nil {
		body["stop"] = req.Stop
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := strings.TrimRight(provider.BaseURL, "/") + "/chat/completions"
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	return httpReq, nil
}

func (a *QwenAdapter) ConvertResponse(resp *http.Response, requestModel string) (*model.ChatResponse, *model.ChatUsage, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("qwen api error (status %d): %s", resp.StatusCode, string(body))
	}

	var result model.ChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, nil, err
	}

	result.Model = requestModel
	usage := result.Usage

	return &result, usage, nil
}

func (a *QwenAdapter) ConvertStreamResponse(resp *http.Response, writer http.ResponseWriter, flusher http.Flusher, requestModel string) (*model.ChatUsage, error) {
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("qwen stream error (status %d): %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	var usage model.ChatUsage

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			fmt.Fprintf(writer, "data: [DONE]\n\n")
			flusher.Flush()
			break
		}

		var chunk map[string]interface{}
		if err := json.Unmarshal([]byte(data), &chunk); err == nil {
			chunk["model"] = requestModel
			if u, ok := chunk["usage"]; ok && u != nil {
				if um, ok := u.(map[string]interface{}); ok {
					if v, ok := um["prompt_tokens"].(float64); ok {
						usage.PromptTokens = int(v)
					}
					if v, ok := um["completion_tokens"].(float64); ok {
						usage.CompletionTokens = int(v)
					}
					if v, ok := um["total_tokens"].(float64); ok {
						usage.TotalTokens = int(v)
					}
				}
			}
			modified, _ := json.Marshal(chunk)
			fmt.Fprintf(writer, "data: %s\n\n", string(modified))
		} else {
			fmt.Fprintf(writer, "data: %s\n\n", data)
		}
		flusher.Flush()
	}

	return &usage, scanner.Err()
}
