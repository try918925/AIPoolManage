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
	"time"

	"github.com/google/uuid"
)

type OpenAIAdapter struct{}

func (a *OpenAIAdapter) ProviderType() string { return "openai" }

func (a *OpenAIAdapter) ConvertRequest(req *model.ChatRequest, channel *model.ProviderModel, provider *model.Provider, apiKey string) (*http.Request, error) {
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
	if req.N != nil {
		body["n"] = *req.N
	}
	if req.Stop != nil {
		body["stop"] = req.Stop
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := strings.TrimRight(provider.BaseURL, "/") + "/v1/chat/completions"
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	if provider.OrgID != "" {
		httpReq.Header.Set("OpenAI-Organization", provider.OrgID)
	}

	return httpReq, nil
}

func (a *OpenAIAdapter) ConvertResponse(resp *http.Response, requestModel string) (*model.ChatResponse, *model.ChatUsage, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("openai api error (status %d): %s", resp.StatusCode, string(body))
	}

	var result model.ChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, nil, err
	}

	result.Model = requestModel
	usage := result.Usage

	return &result, usage, nil
}

func (a *OpenAIAdapter) ConvertStreamResponse(resp *http.Response, writer http.ResponseWriter, flusher http.Flusher, requestModel string) (*model.ChatUsage, error) {
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai stream error (status %d): %s", resp.StatusCode, string(body))
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

		// Parse to extract usage if present, then re-emit with correct model
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

// Helper to generate an OpenAI-style chat completion ID
func generateChatID() string {
	return "chatcmpl-" + uuid.New().String()[:12]
}

func nowUnix() int64 {
	return time.Now().Unix()
}
