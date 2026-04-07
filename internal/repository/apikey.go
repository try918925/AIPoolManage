package repository

import (
	"awesomeProject/internal/model"
	"time"

	"gorm.io/gorm"
)

type APIKeyRepo struct {
	db *gorm.DB
}

func NewAPIKeyRepo(db *gorm.DB) *APIKeyRepo {
	return &APIKeyRepo{db: db}
}

func (r *APIKeyRepo) Create(key *model.UserAPIKey) error {
	return r.db.Create(key).Error
}

func (r *APIKeyRepo) FindByHash(hash string) (*model.UserAPIKey, error) {
	var key model.UserAPIKey
	err := r.db.Where("key_hash = ?", hash).First(&key).Error
	return &key, err
}

func (r *APIKeyRepo) FindByUserID(userID string) ([]model.UserAPIKey, error) {
	var keys []model.UserAPIKey
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&keys).Error
	return keys, err
}

func (r *APIKeyRepo) FindAll() ([]model.UserAPIKey, error) {
	var keys []model.UserAPIKey
	err := r.db.Order("created_at DESC").Find(&keys).Error
	return keys, err
}

func (r *APIKeyRepo) FindByID(id int64) (*model.UserAPIKey, error) {
	var key model.UserAPIKey
	err := r.db.First(&key, id).Error
	return &key, err
}

func (r *APIKeyRepo) Delete(id int64, userID string) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.UserAPIKey{}).Error
}

func (r *APIKeyRepo) AdminDelete(id int64) error {
	return r.db.Delete(&model.UserAPIKey{}, id).Error
}

func (r *APIKeyRepo) UpdateLastUsed(id int64) {
	now := time.Now()
	r.db.Model(&model.UserAPIKey{}).Where("id = ?", id).Update("last_used_at", now)
}

func (r *APIKeyRepo) IncrementQuotaUsed(id int64, tokens int64) error {
	return r.db.Model(&model.UserAPIKey{}).
		Where("id = ?", id).
		UpdateColumn("quota_used", gorm.Expr("quota_used + ?", tokens)).Error
}
