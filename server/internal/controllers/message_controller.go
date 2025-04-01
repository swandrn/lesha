package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"lesha.com/server/internal/entity"
	"lesha.com/server/internal/services"
)

type MessageController struct {
	messageService *services.MessageService
}

func NewMessageController(messageService *services.MessageService) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}

// GetChannelMessages returns all messages in a channel
func (c *MessageController) GetChannelMessages(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	channelId := vars["channelID"]

	messages, err := c.messageService.GetChannelMessages(channelId)
	if err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

// GetMessage returns a specific message by ID
func (c *MessageController) GetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["id"]

	message, err := c.messageService.GetMessage(messageId)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}

// CreateMessage creates a new message
func (c *MessageController) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message entity.Message
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	user, err := services.ExtractUserFromToken(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusUnauthorized)
		return
	}

	userID := user.ID
	message.UserID = userID
	channelID := r.FormValue("channelID")
	channelIDUint, err := strconv.ParseUint(channelID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid channel ID", http.StatusBadRequest)
		return
	}
	message.ChannelID = uint(channelIDUint)
	message.Content = r.FormValue("content")

	if err := c.messageService.CreateMessage(&message); err != nil {
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

// PinMessage pins a message
func (c *MessageController) PinMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["id"]

	message, err := c.messageService.GetMessage(messageId)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if err := c.messageService.PinMessage(message); err != nil {
		http.Error(w, "Failed to pin message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Message pinned successfully",
	})
}

// UnpinMessage unpins a message
func (c *MessageController) UnpinMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["id"]

	message, err := c.messageService.GetMessage(messageId)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if err := c.messageService.UnpinMessage(message); err != nil {
		http.Error(w, "Failed to unpin message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Message unpinned successfully",
	})
}

// AddReaction adds a reaction to a message
func (c *MessageController) AddReaction(w http.ResponseWriter, r *http.Request) {
	var reaction entity.Reaction

	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	user, err := services.ExtractUserFromToken(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	messageId := vars["id"]
	messageIdUint, err := strconv.ParseUint(messageId, 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	userId := user.ID
	reaction.UserID = userId
	reaction.Emoji = r.FormValue("emoji")
	reaction.MessageID = uint(messageIdUint)

	if err := c.messageService.AddReaction(&reaction); err != nil {
		http.Error(w, "Failed to add reaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reaction)
}

// RemoveReaction removes a reaction from a message
func (c *MessageController) RemoveReaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["id"]
	reactionId := vars["reactionId"]

	reactions, err := c.messageService.GetReactions(messageId)
	if err != nil {
		http.Error(w, "Reaction not found", http.StatusNotFound)
		return
	}

	// Find the specific reaction by ID
	var reactionToRemove *entity.Reaction
	for _, r := range reactions {
		if strconv.FormatUint(uint64(r.ID), 10) == reactionId {
			reactionToRemove = &r
			break
		}
	}

	if reactionToRemove == nil {
		http.Error(w, "Reaction not found", http.StatusNotFound)
		return
	}

	if err := c.messageService.RemoveReaction(reactionToRemove); err != nil {
		http.Error(w, "Failed to remove reaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Reaction removed successfully",
	})
}

// AddMedia adds a media to a message
func (c *MessageController) AddMedia(w http.ResponseWriter, r *http.Request) {
	var media entity.Media

	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	user, err := services.ExtractUserFromToken(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	messageId := vars["id"]

	message, err := c.messageService.GetMessage(messageId)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if message.UserID != user.ID {
		http.Error(w, "You are not allowed to add media to this message", http.StatusForbidden)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create upload directory for media
	uploadDir := "uploads/messages"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
	filepath := path.Join(uploadDir, filename)

	// Create the file
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the created file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Get file extension and determine media type
	filename = handler.Filename
	var ext string
	for _, e := range []string{".jpg", ".jpeg", ".png", ".gif", ".mp4", ".webm", ".mp3", ".wav"} {
		if strings.HasSuffix(strings.ToLower(filename), e) {
			ext = e
			break
		}
	}
	if ext == "" {
		http.Error(w, "Invalid file extension", http.StatusBadRequest)
		return
	}

	// Determine media type based on extension
	var mediaType string
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		mediaType = "image"
	case ".mp4", ".webm":
		mediaType = "video"
	case ".mp3", ".wav":
		mediaType = "audio"
	default:
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	media.MessageID = message.ID
	media.Type = mediaType
	media.Extension = ext[1:] // Remove the dot from extension
	media.Url = filepath

	if err := c.messageService.AddMedia(&media); err != nil {
		http.Error(w, "Failed to add media", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(media)
}
