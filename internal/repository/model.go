package repository

import (
	"awesomeProject/internal/model"

	"gorm.io/gorm"
)

type ModelRepo struct {
	db *gorm.DB
}

func NewModelRepo(db *gorm.DB) *ModelRepo {
	return &ModelRepo{db: db}
}

func (r *ModelRepo) Create(m *model.ProviderModel) error {
	return r.db.Create(m).Error
}

func (r *ModelRepo) FindByID(id int64) (*model.ProviderModel, error) {
	var m model.ProviderModel
	err := r.db.Preload("Provider").First(&m, id).Error
	return &m, err
}

func (r *ModelRepo) FindByProviderID(providerID int64) ([]model.ProviderModel, error) {
	var models []model.ProviderModel
	err := r.db.Where("provider_id = ?", providerID).Order("id ASC").Find(&models).Error
	return models, err
}

func (r *ModelRepo) FindAllByName(modelName string) ([]model.ProviderModel, error) {
	var models []model.ProviderModel
	err := r.db.Preload("Provider").
		Where("model_name = ? AND enabled = true", modelName).
		Order("priority ASC, weight DESC").
		Find(&models).Error
	return models, err
}

func (r *ModelRepo) FindDistinctModelNames() ([]string, error) {
	var names []string
	err := r.db.Model(&model.ProviderModel{}).
		Where("enabled = true").
		Distinct("model_name").
		Pluck("model_name", &names).Error
	return names, err
}

func (r *ModelRepo) FindAllEnabled() ([]model.ProviderModel, error) {
	var models []model.ProviderModel
	err := r.db.Preload("Provider").Where("enabled = true").Find(&models).Error
	return models, err
}

func (r *ModelRepo) Update(m *model.ProviderModel) error {
	return r.db.Save(m).Error
}

func (r *ModelRepo) ExistsByName(modelName string) bool {
	var count int64
	r.db.Model(&model.ProviderModel{}).Where("model_name = ?", modelName).Count(&count)
	return count > 0
}

func (r *ModelRepo) ExistsByNameExcludeID(modelName string, excludeID int64) bool {
	var count int64
	r.db.Model(&model.ProviderModel{}).Where("model_name = ? AND id != ?", modelName, excludeID).Count(&count)
	return count > 0
}

func (r *ModelRepo) Delete(id int64) error {
	return r.db.Delete(&model.ProviderModel{}, id).Error
}
