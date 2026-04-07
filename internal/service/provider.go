package service

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/pkg/crypto"
	"awesomeProject/internal/repository"
)

type ProviderService struct {
	repo      *repository.ProviderRepo
	masterKey []byte
}

func NewProviderService(repo *repository.ProviderRepo, masterKey []byte) *ProviderService {
	return &ProviderService{repo: repo, masterKey: masterKey}
}

type CreateProviderReq struct {
	Name    string                 `json:"name" binding:"required"`
	Type    string                 `json:"type" binding:"required"`
	BaseURL string                 `json:"base_url" binding:"required"`
	APIKey  string                 `json:"api_key" binding:"required"`
	OrgID   string                 `json:"org_id"`
	Config  map[string]interface{} `json:"config"`
}

type UpdateProviderReq struct {
	Name    *string                `json:"name"`
	Type    *string                `json:"type"`
	BaseURL *string                `json:"base_url"`
	APIKey  *string                `json:"api_key"`
	OrgID   *string                `json:"org_id"`
	Enabled *bool                  `json:"enabled"`
	Config  map[string]interface{} `json:"config"`
}

func (s *ProviderService) Create(req *CreateProviderReq) (*model.Provider, error) {
	encrypted, err := crypto.Encrypt(req.APIKey, s.masterKey)
	if err != nil {
		return nil, err
	}

	provider := &model.Provider{
		Name:            req.Name,
		Type:            req.Type,
		BaseURL:         req.BaseURL,
		APIKeyEncrypted: encrypted,
		OrgID:           req.OrgID,
		Enabled:         true,
	}

	if req.Config != nil {
		configJSON, _ := jsonMarshal(req.Config)
		provider.Config = configJSON
	}

	if err := s.repo.Create(provider); err != nil {
		return nil, err
	}

	return provider, nil
}

func (s *ProviderService) List() ([]model.Provider, error) {
	return s.repo.FindAll()
}

func (s *ProviderService) GetByID(id int64) (*model.Provider, error) {
	return s.repo.FindByID(id)
}

func (s *ProviderService) Update(id int64, req *UpdateProviderReq) (*model.Provider, error) {
	provider, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		provider.Name = *req.Name
	}
	if req.Type != nil {
		provider.Type = *req.Type
	}
	if req.BaseURL != nil {
		provider.BaseURL = *req.BaseURL
	}
	if req.OrgID != nil {
		provider.OrgID = *req.OrgID
	}
	if req.Enabled != nil {
		provider.Enabled = *req.Enabled
	}
	if req.APIKey != nil {
		encrypted, err := crypto.Encrypt(*req.APIKey, s.masterKey)
		if err != nil {
			return nil, err
		}
		provider.APIKeyEncrypted = encrypted
	}
	if req.Config != nil {
		configJSON, _ := jsonMarshal(req.Config)
		provider.Config = configJSON
	}

	if err := s.repo.Update(provider); err != nil {
		return nil, err
	}

	return provider, nil
}

func (s *ProviderService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *ProviderService) DecryptAPIKey(provider *model.Provider) (string, error) {
	return crypto.Decrypt(provider.APIKeyEncrypted, s.masterKey)
}
