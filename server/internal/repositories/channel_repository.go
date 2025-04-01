// repositories/channel_repository.go
package repositories

import (
	"gorm.io/gorm"
	"lesha.com/server/internal/entity"
)

type ChannelRepository struct {
	DB *gorm.DB
}

func NewChannelRepository(db *gorm.DB) *ChannelRepository {
	return &ChannelRepository{DB: db}
}

func (repo *ChannelRepository) CreateChannel(channel *entity.Channel) error {
	return repo.DB.Create(channel).Error
}
func (repo *ChannelRepository) GetChannel(channelId string) (*entity.Channel, error) {
	var channel entity.Channel
	err := repo.DB.Where("id = ?", channelId).First(&channel).Error
	if err != nil {
		return nil, err
	}
	return &channel, nil
}
func (repo *ChannelRepository) GetChannels() ([]entity.Channel, error) {
	var channels []entity.Channel
	err := repo.DB.Find(&channels).Error
	if err != nil {
		return nil, err
	}
	return channels, nil
}
func (repo *ChannelRepository) UpdateChannel(channel *entity.Channel) error {
	return repo.DB.Save(channel).Error
}
func (repo *ChannelRepository) DeleteChannel(channel *entity.Channel) error {
	return repo.DB.Delete(channel).Error
}
func (repo *ChannelRepository) AddUserToChannel(channelID uint, userID uint) error {
	return repo.DB.Exec("INSERT INTO user_channels (channel_id, user_id) VALUES (?, ?)", channelID, userID).Error
}
