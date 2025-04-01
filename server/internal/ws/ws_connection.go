package ws

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"lesha.com/server/internal/entity"
	"lesha.com/server/internal/services"
)

var ChannelClients = make(map[string]map[*Client]bool)

type Client struct {
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   uint
	Channels map[string]bool
}

func (c *Client) readPump(db *gorm.DB) {
	defer func() { c.Conn.Close() }()
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		c.handleMessage(db, msg)
	}
}

func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("write error:", err)
				return
			}
		}
	}
}

func (c *Client) joinChannel(channelName string) {
	if _, ok := ChannelClients[channelName]; !ok {
		ChannelClients[channelName] = make(map[*Client]bool)
	}
	ChannelClients[channelName][c] = true
	c.Channels[channelName] = true

	log.Printf("User %d joined channel %s", c.UserID, channelName)
}

func HandleWebSocket(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := Upgrade(w, r)
		if err != nil {
			http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value
		user, err := services.ExtractUserFromToken(tokenString)
		if err != nil {
			return
		}

		client := &Client{
			Conn:     conn,
			Send:     make(chan []byte, 256),
			UserID:   user.ID,
			Channels: make(map[string]bool),
		}

		if err := db.Preload("Servers.Channels").First(&user, user.ID).Error; err != nil {
			log.Println("failed to fetch user servers/channels:", err)
			return
		}

		go client.writePump()
		client.readPump(db)
	}
}

func (c *Client) handleMessage(db *gorm.DB, raw []byte) {
	var incoming struct {
		Type      string `json:"type"`
		ChannelID uint   `json:"channel_id"`
		MessageID uint   `json:"message_id"`
		Content   string `json:"content"`
		File      string `json:"file"`
		Filename  string `json:"filename"`
		Reaction  string `json:"reaction"`
	}

	if err := json.Unmarshal(raw, &incoming); err != nil {
		log.Println("Invalid message:", err)
		return
	}

	fmt.Println(incoming)

	switch incoming.Type {
	case "MESSAGE":
		messageService := services.NewMessageService(db)
		message := entity.Message{
			UserID:    c.UserID,
			ChannelID: incoming.ChannelID,
			Content:   incoming.Content,
			Pinned:    false,
		}

		if err := messageService.CreateMessage(&message); err != nil {
			log.Println("Failed to save message:", err)
			return
		}

		if incoming.File != "" && incoming.Filename != "" {
			parts := strings.SplitN(incoming.File, ",", 2)
			if len(parts) != 2 {
				log.Println("Invalid base64 data")
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				log.Println("Failed to decode base64:", err)
				return
			}

			uploadDir := "uploads/messages"
			if err := os.MkdirAll(uploadDir, 0755); err != nil {
				log.Println("Failed to create upload directory:", err)
				return
			}

			filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), incoming.Filename)
			filePath := filepath.Join(uploadDir, filename)

			file, err := os.Create(filePath)
			if err != nil {
				log.Println("Failed to create file:", err)
				return
			}
			defer file.Close()

			if _, err := file.Write(decoded); err != nil {
				log.Println("Failed to write file:", err)
				return
			}

			var ext string
			for _, e := range []string{".jpg", ".jpeg", ".png", ".gif", ".mp4", ".webm", ".mp3", ".wav"} {
				if strings.HasSuffix(strings.ToLower(incoming.Filename), e) {
					ext = e
					break
				}
			}

			if ext == "" {
				log.Println("Invalid file extension")
				return
			}

			var mediaType string
			switch ext {
			case ".jpg", ".jpeg", ".png", ".gif":
				mediaType = "image"
			case ".mp4", ".webm":
				mediaType = "video"
			case ".mp3", ".wav":
				mediaType = "audio"
			default:
				log.Println("Unsupported file type")
				return
			}

			media := entity.Media{
				MessageID: message.ID,
				Type:      mediaType,
				Extension: ext[1:],
				Url:       filePath,
			}

			if err := messageService.AddMedia(&media); err != nil {
				log.Println("Failed to save media:", err)
				return
			}
		}

		updatedMessage, err := messageService.GetMessage(fmt.Sprintf("%d", message.ID))
		if err != nil {
			log.Println("Failed to get updated message:", err)
		}

		messageResponse := updatedMessage.ToResponse()

		payload, _ := json.Marshal(struct {
			Type      string                 `json:"type"`
			ID        uint                   `json:"id"`
			ChannelID uint                   `json:"channel_id"`
			SenderID  uint                   `json:"sender"`
			User      entity.UserResponse    `json:"user"`
			Content   string                 `json:"content"`
			Timestamp time.Time              `json:"timestamp"`
			Medias    []entity.MediaResponse `json:"medias"`
		}{
			Type:      "MESSAGE",
			ID:        messageResponse.ID,
			ChannelID: messageResponse.ChannelID,
			SenderID:  messageResponse.User.ID,
			User:      messageResponse.User,
			Content:   messageResponse.Content,
			Timestamp: messageResponse.CreatedAt,
			Medias:    messageResponse.Medias,
		})

		broadcastToChannel(db, message.ChannelID, payload)

	case "JOIN_CHANNEL":
		channelService := services.NewChannelService(db)

		channel, err := channelService.GetChannel(incoming.ChannelID)
		if err != nil {
			log.Println("Failed to find channel:", err)
			return
		}
		c.joinChannel(channel.Name)

	case "REACTION":
		if incoming.MessageID == 0 || incoming.Reaction == "" {
			log.Println("Invalid reaction data")
			return
		}

		messageService := services.NewMessageService(db)

		// Create a new reaction
		reaction := entity.Reaction{
			UserID:    c.UserID,
			MessageID: incoming.MessageID,
			Emoji:     incoming.Reaction,
		}

		// Add the reaction to the database
		if err := messageService.AddReaction(&reaction); err != nil {
			log.Println("Failed to add reaction:", err)
			return
		}

		// Get the updated message with the new reaction
		messageId := fmt.Sprintf("%d", incoming.MessageID)
		message, err := messageService.GetMessage(messageId)
		if err != nil {
			log.Println("Failed to get updated message:", err)
			return
		}

		// Convert to response format
		messageResponse := message.ToResponse()

		// Broadcast the updated message to all clients in the channel
		payload, _ := json.Marshal(struct {
			Type      string                    `json:"type"`
			ID        uint                      `json:"id"`
			ChannelID uint                      `json:"channel_id"`
			SenderID  uint                      `json:"sender"`
			User      entity.UserResponse       `json:"user"`
			Content   string                    `json:"content"`
			Timestamp time.Time                 `json:"timestamp"`
			Medias    []entity.MediaResponse    `json:"medias"`
			Reactions []entity.ReactionResponse `json:"reactions"`
		}{
			Type:      "MESSAGE_UPDATE",
			ID:        messageResponse.ID,
			ChannelID: messageResponse.ChannelID,
			SenderID:  messageResponse.User.ID,
			User:      messageResponse.User,
			Content:   messageResponse.Content,
			Timestamp: messageResponse.CreatedAt,
			Medias:    messageResponse.Medias,
			Reactions: messageResponse.Reactions,
		})

		broadcastToChannel(db, messageResponse.ChannelID, payload)
	}
}

func broadcastToChannel(db *gorm.DB, channelId uint, message []byte) {
	channelService := services.NewChannelService(db)
	channel, err := channelService.GetChannel(channelId)
	if err != nil {
		return
	}

	clients := ChannelClients[channel.Name]
	for client := range clients {
		select {
		case client.Send <- message:
		default:
			log.Printf("Client %d buffer full", client.UserID)
		}
	}
}
