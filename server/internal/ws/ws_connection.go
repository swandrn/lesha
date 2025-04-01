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
		Content   string `json:"content"`
		File      string `json:"file"`
		Filename  string `json:"filename"`
	}

	if err := json.Unmarshal(raw, &incoming); err != nil {
		log.Println("Invalid message:", err)
		return
	}

	switch incoming.Type {
	case "MESSAGE":
		messageService := services.NewMessageService(db)
		message := entity.Message{
			UserID:    c.UserID,
			ChannelID: incoming.ChannelID,
			Content:   incoming.Content,
			Pinned:    false,
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

			uploadDir := "uploads/servers"
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

			message.Content = filename
		}

		if err := messageService.CreateMessage(&message); err != nil {
			log.Println("Failed to save message:", err)
			return
		}

		payload, _ := json.Marshal(struct {
			Type      string    `json:"type"`
			ChannelID uint      `json:"channel_id"`
			SenderID  uint      `json:"sender"`
			Content   string    `json:"content"`
			Timestamp time.Time `json:"timestamp"`
		}{
			Type:      "MESSAGE",
			ChannelID: message.ChannelID,
			SenderID:  message.UserID,
			Content:   message.Content,
			Timestamp: message.CreatedAt,
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
