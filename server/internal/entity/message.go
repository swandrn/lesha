package entity

import "time"

// MessageResponse represents the cleaned up message response
type MessageResponse struct {
	ID        uint               `json:"id"`
	CreatedAt time.Time          `json:"createdAt"`
	User      UserResponse       `json:"user"`
	Reactions []ReactionResponse `json:"reactions"`
	Medias    []MediaResponse    `json:"medias"`
	ChannelID uint               `json:"channelId"`
	Content   string             `json:"content"`
	Pinned    bool               `json:"pinned"`
}

// UserResponse represents the cleaned up user response
type UserResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

// ReactionResponse represents the cleaned up reaction response
type ReactionResponse struct {
	ID     uint   `json:"id"`
	Emoji  string `json:"emoji"`
	UserID uint   `json:"userId"`
}

// MediaResponse represents the cleaned up media response
type MediaResponse struct {
	ID        uint   `json:"id"`
	Type      string `json:"type"`
	Extension string `json:"extension"`
	Url       string `json:"url"`
}

// ToResponse converts a Message to MessageResponse
func (m *Message) ToResponse() MessageResponse {
	reactions := make([]ReactionResponse, len(m.Reactions))
	for i, r := range m.Reactions {
		reactions[i] = ReactionResponse{
			ID:     r.ID,
			Emoji:  r.Emoji,
			UserID: r.UserID,
		}
	}

	medias := make([]MediaResponse, len(m.Medias))
	for i, media := range m.Medias {
		medias[i] = MediaResponse{
			ID:        media.ID,
			Type:      media.Type,
			Extension: media.Extension,
			Url:       media.Url,
		}
	}

	return MessageResponse{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		User: UserResponse{
			ID:          m.User.ID,
			Name:        m.User.Name,
			DisplayName: m.User.DisplayName,
		},
		Reactions: reactions,
		Medias:    medias,
		ChannelID: m.ChannelID,
		Content:   m.Content,
		Pinned:    m.Pinned,
	}
}
