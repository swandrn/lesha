package services

import (
	"gorm.io/gorm"
	"lesha.com/server/internal/entity"
	"lesha.com/server/internal/repositories"
)

type MessageService struct {
	DB *gorm.DB
}

func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{DB: db}
}

func (service *MessageService) CreateMessage(message *entity.Message) error {
	messageRepository := repositories.NewMessageRepository(service.DB)
	return messageRepository.CreateMessage(message)
}

func (service *MessageService) GetMessage(messageId string) (*entity.Message, error) {
	messageRepository := repositories.NewMessageRepository(service.DB)
	return messageRepository.GetMessage(messageId)
}

func (service *MessageService) GetChannelMessages(channelId string) ([]entity.Message, error) {
	messageRepository := repositories.NewMessageRepository(service.DB)
	return messageRepository.GetChannelMessages(channelId)
}

func (service *MessageService) PinMessage(message *entity.Message) error {
	messageRepository := repositories.NewMessageRepository(service.DB)
	return messageRepository.PinMessage(message)
}

func (service *MessageService) UnpinMessage(message *entity.Message) error {
	messageRepository := repositories.NewMessageRepository(service.DB)
	return messageRepository.UnpinMessage(message)
}

func (service *MessageService) AddReaction(reaction *entity.Reaction) error {
	messageRepository := repositories.NewMessageRepository(service.DB)
	return messageRepository.AddReaction(reaction)
}

func (service *MessageService) RemoveReaction(reaction *entity.Reaction) error {
	messageRepository := repositories.NewMessageRepository(service.DB)
	return messageRepository.RemoveReaction(reaction)
}

func (service *MessageService) GetReactions(messageId string) ([]entity.Reaction, error) {
	messageRepository := repositories.NewMessageRepository(service.DB)
	return messageRepository.GetReactions(messageId)
}

func (service *MessageService) AddMedia(media *entity.Media) error {
	messageRepository := repositories.NewMessageRepository(service.DB)
	return messageRepository.AddMedia(media)
}

func (service *MessageService) RemoveMedia(media *entity.Media) error {
	messageRepository := repositories.NewMessageRepository(service.DB)
	return messageRepository.RemoveMedia(media)
}

func (service *MessageService) GetMedia(mediaId string) (*entity.Media, error) {
	messageRepository := repositories.NewMessageRepository(service.DB)
	return messageRepository.GetMedia(mediaId)
}
