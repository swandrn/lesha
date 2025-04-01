package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"lesha.com/server/internal/database"
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

func (c *Client) readPump() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		c.handleMessage(msg)
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
				// Channel closed
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

// Function called when a user upgrades from http to ws
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrade(w, r)
	if err != nil {
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}

	// 🔐 Extract user ID from auth (replace with real logic)
	userID := uint(1)

	// Create client
	client := &Client{
		Conn:     conn,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		Channels: make(map[string]bool),
	}

	// Auto-join all server channels
	db := database.Connect()
	var user entity.User
	if err := db.Preload("Servers.Channels").First(&user, userID).Error; err != nil {
		log.Println("failed to fetch user servers/channels:", err)
		return
	}

	for _, server := range user.Servers {
		for _, channel := range server.Channels {
			client.joinChannel(channel.Name)
		}
	}

	// Start reading/writing
	go client.writePump()
	client.readPump()
}

func (c *Client) handleMessage(raw []byte) {
	// Receive a message from client
	var incoming struct {
		Type      string `json:"type"`
		ChannelID uint   `json:"channel_id"`
		Content   string `json:"content"`
	}

	if err := json.Unmarshal(raw, &incoming); err != nil {
		log.Println("Invalid message:", err)
		return
	}

	switch incoming.Type {
	case "MESSAGE":
		db := database.Connect()
		messageService := services.NewMessageService(db)

		message := entity.Message{
			UserID:    c.UserID,
			Reactions: []entity.Reaction{},
			Medias:    []entity.Media{},
			ChannelID: incoming.ChannelID,
			Content:   incoming.Content,
			Pinned:    false,
		}

		// Save to DB
		if err := messageService.CreateMessage(&message); err != nil {
			log.Println("Failed to save message:", err)
			return
		}

		// Prepare broadcast payload
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

		broadcastToChannel(message.ChannelID, payload)
	}
}

func broadcastToChannel(channelId uint, message []byte) {
	db := database.Connect()
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
