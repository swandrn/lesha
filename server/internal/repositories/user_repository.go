// repositories/user_repository.go
package repositories

import (
	"gorm.io/gorm"
	"lesha.com/server/internal/entity"
)

type BlacklistedTokenRepository struct {
	DB *gorm.DB
}

func NewBlacklistedTokenRepository(db *gorm.DB) *BlacklistedTokenRepository {
	return &BlacklistedTokenRepository{DB: db}
}

func (repo *BlacklistedTokenRepository) GetBlacklistedToken(token string) (*entity.BlacklistedToken, error) {
	var blacklistedToken entity.BlacklistedToken
	if err := repo.DB.Where("token = ?", token).First(&blacklistedToken).Error; err != nil {
		return nil, err
	}
	return &blacklistedToken, nil
}
func (repo *BlacklistedTokenRepository) GetAllBlacklistedTokens() ([]entity.BlacklistedToken, error) {
	var blacklistedTokens []entity.BlacklistedToken
	if err := repo.DB.Find(&blacklistedTokens).Error; err != nil {
		return nil, err
	}
	return blacklistedTokens, nil
}
func (repo *BlacklistedTokenRepository) CreateBlacklistedToken(token string) error {
	return repo.DB.Create(&entity.BlacklistedToken{Token: token}).Error
}
