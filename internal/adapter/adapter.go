package adapter

import (
	"awesomeProject/internal/model"
	"net/http"
)

type Adapter interface {
	ConvertRequest(req *model.ChatRequest, channel *model.ProviderModel, provider *model.Provider, apiKey string) (*http.Request, error)
	ConvertResponse(resp *http.Response, requestModel string) (*model.ChatResponse, *model.ChatUsage, error)
	ConvertStreamResponse(resp *http.Response, writer http.ResponseWriter, flusher http.Flusher, requestModel string) (*model.ChatUsage, error)
	ProviderType() string
}
