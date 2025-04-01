package services

import (
	"gorm.io/gorm"
	"lesha.com/server/internal/entity"
	"lesha.com/server/internal/repositories"
)

type ChannelService struct {
	DB *gorm.DB
}

func NewChannelService(db *gorm.DB) *ChannelService {
	return &ChannelService{DB: db}
}

func (service *ChannelService) CreateChannel(channel *entity.Channel) error {
	channelRepository := repositories.NewChannelRepository(service.DB)
	return channelRepository.CreateChannel(channel)
}

func (service *ChannelService) GetChannel(channelId uint) (*entity.Channel, error) {
	channelRepository := repositories.NewChannelRepository(service.DB)
	return channelRepository.GetChannel(channelId)
}

func (service *ChannelService) GetChannels() ([]entity.Channel, error) {
	channelRepository := repositories.NewChannelRepository(service.DB)
	return channelRepository.GetChannels()
}

func (service *ChannelService) UpdateChannel(channel *entity.Channel) error {
	channelRepository := repositories.NewChannelRepository(service.DB)
	return channelRepository.UpdateChannel(channel)
}

func (service *ChannelService) DeleteChannel(channel *entity.Channel) error {
	channelRepository := repositories.NewChannelRepository(service.DB)
	return channelRepository.DeleteChannel(channel)
}

func (service *ChannelService) AddUserToChannel(channelID uint, userID uint) error {
	channelRepository := repositories.NewChannelRepository(service.DB)
	return channelRepository.AddUserToChannel(channelID, userID)
}
