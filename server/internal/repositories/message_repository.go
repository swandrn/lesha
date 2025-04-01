// repositories/message_repository.go
package repositories

import (
	"gorm.io/gorm"
	"lesha.com/server/internal/entity"
)

type MessageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

func (repo *MessageRepository) CreateMessage(message *entity.Message) error {
	return repo.DB.Create(message).Error
}
func (repo *MessageRepository) UpdateMessage(message *entity.Message) error {
	return repo.DB.Save(message).Error
}
func (repo *MessageRepository) DeleteMessage(message *entity.Message) error {
	return repo.DB.Delete(message).Error
}
func (repo *MessageRepository) GetMessage(messageId string) (*entity.Message, error) {
	var message entity.Message
	err := repo.DB.Where("id = ?", messageId).Preload("Medias").Preload("Reactions").Preload("User").First(&message).Error
	return &message, err
}
func (repo *MessageRepository) GetChannelMessages(channelId string) ([]entity.Message, error) {
	var messages []entity.Message

	err := repo.DB.Where("channel_id = ?", channelId).Preload("Medias").Preload("Reactions").Preload("User").Find(&messages).Error
	return messages, err
}

// Pin message
func (repo *MessageRepository) PinMessage(message *entity.Message) error {
	return repo.DB.Model(message).Update("pinned", true).Error
}
func (repo *MessageRepository) UnpinMessage(message *entity.Message) error {
	return repo.DB.Model(message).Update("pinned", false).Error
}

// Reactions
func (repo *MessageRepository) AddReaction(reaction *entity.Reaction) error {
	return repo.DB.Create(reaction).Error
}
func (repo *MessageRepository) RemoveReaction(reaction *entity.Reaction) error {
	return repo.DB.Delete(reaction).Error
}
func (repo *MessageRepository) GetReactions(messageId string) ([]entity.Reaction, error) {
	var reactions []entity.Reaction
	err := repo.DB.Where("message_id = ?", messageId).Find(&reactions).Error
	return reactions, err
}

// Media
func (repo *MessageRepository) AddMedia(media *entity.Media) error {
	return repo.DB.Create(media).Error
}
func (repo *MessageRepository) RemoveMedia(media *entity.Media) error {
	return repo.DB.Delete(media).Error
}
func (repo *MessageRepository) GetMedia(mediaId string) (*entity.Media, error) {
	var media entity.Media
	err := repo.DB.Where("id = ?", mediaId).First(&media).Error
	return &media, err
}
