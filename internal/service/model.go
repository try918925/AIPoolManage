package service

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"fmt"
)

type ModelService struct {
	modelRepo    *repository.ModelRepo
	providerRepo *repository.ProviderRepo
}

func NewModelService(modelRepo *repository.ModelRepo, providerRepo *repository.ProviderRepo) *ModelService {
	return &ModelService{modelRepo: modelRepo, providerRepo: providerRepo}
}

type CreateModelReq struct {
	ModelName        string                 `json:"model_name" binding:"required"`
	ModelID          string                 `json:"model_id" binding:"required"`
	ModelType        string                 `json:"model_type"`
	MaxContextTokens *int                   `json:"max_context_tokens"`
	InputPrice       *float64               `json:"input_price"`
	OutputPrice      *float64               `json:"output_price"`
	Weight           int                    `json:"weight"`
	Priority         int                    `json:"priority"`
	Config           map[string]interface{} `json:"config"`
}

type UpdateModelReq struct {
	ModelName        *string                `json:"model_name"`
	ModelID          *string                `json:"model_id"`
	ModelType        *string                `json:"model_type"`
	Enabled          *bool                  `json:"enabled"`
	Weight           *int                   `json:"weight"`
	Priority         *int                   `json:"priority"`
	MaxContextTokens *int                   `json:"max_context_tokens"`
	InputPrice       *float64               `json:"input_price"`
	OutputPrice      *float64               `json:"output_price"`
	Config           map[string]interface{} `json:"config"`
}

func (s *ModelService) Create(providerID int64, req *CreateModelReq) (*model.ProviderModel, error) {
	if s.modelRepo.ExistsByName(req.ModelName) {
		return nil, fmt.Errorf("模型名称 '%s' 已存在", req.ModelName)
	}

	mt := "chat"
	if req.ModelType != "" {
		mt = req.ModelType
	}
	weight := 1
	if req.Weight > 0 {
		weight = req.Weight
	}

	m := &model.ProviderModel{
		ProviderID:       providerID,
		ModelName:        req.ModelName,
		ModelID:          req.ModelID,
		ModelType:        mt,
		Enabled:          true,
		Weight:           weight,
		Priority:         req.Priority,
		MaxContextTokens: req.MaxContextTokens,
		InputPrice:       req.InputPrice,
		OutputPrice:      req.OutputPrice,
	}

	if req.Config != nil {
		configJSON, _ := jsonMarshal(req.Config)
		m.Config = configJSON
	}

	if err := s.modelRepo.Create(m); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *ModelService) ListByProvider(providerID int64) ([]model.ProviderModel, error) {
	return s.modelRepo.FindByProviderID(providerID)
}

func (s *ModelService) Update(id int64, req *UpdateModelReq) (*model.ProviderModel, error) {
	m, err := s.modelRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.ModelName != nil && *req.ModelName != m.ModelName {
		if s.modelRepo.ExistsByNameExcludeID(*req.ModelName, id) {
			return nil, fmt.Errorf("模型名称 '%s' 已存在", *req.ModelName)
		}
		m.ModelName = *req.ModelName
	}
	if req.ModelID != nil {
		m.ModelID = *req.ModelID
	}
	if req.ModelType != nil {
		m.ModelType = *req.ModelType
	}
	if req.Enabled != nil {
		m.Enabled = *req.Enabled
	}
	if req.Weight != nil {
		m.Weight = *req.Weight
	}
	if req.Priority != nil {
		m.Priority = *req.Priority
	}
	if req.MaxContextTokens != nil {
		m.MaxContextTokens = req.MaxContextTokens
	}
	if req.InputPrice != nil {
		m.InputPrice = req.InputPrice
	}
	if req.OutputPrice != nil {
		m.OutputPrice = req.OutputPrice
	}
	if req.Config != nil {
		configJSON, _ := jsonMarshal(req.Config)
		m.Config = configJSON
	}

	if err := s.modelRepo.Update(m); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *ModelService) Delete(id int64) error {
	return s.modelRepo.Delete(id)
}

func (s *ModelService) ListAvailableModels() ([]map[string]interface{}, error) {
	names, err := s.modelRepo.FindDistinctModelNames()
	if err != nil {
		return nil, err
	}

	var models []map[string]interface{}
	for _, name := range names {
		channels, _ := s.modelRepo.FindAllByName(name)
		if len(channels) == 0 {
			continue
		}

		providerType := "unknown"
		if channels[0].Provider != nil {
			providerType = channels[0].Provider.Type
		}

		// Collect unique provider types
		providerTypes := make(map[string]bool)
		for _, ch := range channels {
			if ch.Provider != nil {
				providerTypes[ch.Provider.Type] = true
			}
		}
		providers := make([]string, 0, len(providerTypes))
		for p := range providerTypes {
			providers = append(providers, p)
		}

		models = append(models, map[string]interface{}{
			"id":            name,
			"object":        "model",
			"owned_by":      providerType,
			"created":       channels[0].CreatedAt.Unix(),
			"channel_count": len(channels),
			"providers":     providers,
		})
	}

	return models, nil
}

func (s *ModelService) GetModelDetail(modelName string) (map[string]interface{}, error) {
	channels, err := s.modelRepo.FindAllByName(modelName)
	if err != nil || len(channels) == 0 {
		return nil, err
	}

	providersSet := make(map[string]bool)
	var maxCtx *int
	for _, ch := range channels {
		if ch.Provider != nil {
			providersSet[ch.Provider.Type] = true
		}
		if ch.MaxContextTokens != nil {
			if maxCtx == nil || *ch.MaxContextTokens > *maxCtx {
				maxCtx = ch.MaxContextTokens
			}
		}
	}

	providers := make([]string, 0, len(providersSet))
	for p := range providersSet {
		providers = append(providers, p)
	}

	result := map[string]interface{}{
		"model_name":         modelName,
		"model_type":         channels[0].ModelType,
		"available_channels": len(channels),
		"providers":          providers,
		"status":             "healthy",
	}
	if maxCtx != nil {
		result["max_context_tokens"] = *maxCtx
	}

	return result, nil
}

func (s *ModelService) GetModelChannels(modelName string) ([]model.ProviderModel, error) {
	return s.modelRepo.FindAllByName(modelName)
}
