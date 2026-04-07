package adapter

import "fmt"

type Registry struct {
	adapters map[string]Adapter
}

func NewRegistry() *Registry {
	r := &Registry{adapters: make(map[string]Adapter)}
	r.Register(&OpenAIAdapter{})
	r.Register(&ClaudeAdapter{})
	r.Register(&QwenAdapter{})
	return r
}

func (r *Registry) Register(a Adapter) {
	r.adapters[a.ProviderType()] = a
}

func (r *Registry) Get(providerType string) (Adapter, error) {
	a, ok := r.adapters[providerType]
	if !ok {
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
	return a, nil
}
