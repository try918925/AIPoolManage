package repository

import (
	"awesomeProject/internal/model"

	"gorm.io/gorm"
)

type ProviderRepo struct {
	db *gorm.DB
}

func NewProviderRepo(db *gorm.DB) *ProviderRepo {
	return &ProviderRepo{db: db}
}

func (r *ProviderRepo) Create(p *model.Provider) error {
	return r.db.Create(p).Error
}

func (r *ProviderRepo) FindByID(id int64) (*model.Provider, error) {
	var p model.Provider
	err := r.db.First(&p, id).Error
	return &p, err
}

func (r *ProviderRepo) FindAll() ([]model.Provider, error) {
	var providers []model.Provider
	err := r.db.Order("id ASC").Find(&providers).Error
	return providers, err
}

func (r *ProviderRepo) Update(p *model.Provider) error {
	return r.db.Save(p).Error
}

func (r *ProviderRepo) Delete(id int64) error {
	return r.db.Delete(&model.Provider{}, id).Error
}
