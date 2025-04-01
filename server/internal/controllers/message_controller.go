package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	channelId := vars["channelId"]

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
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

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

// AddReaction adds a reaction to a message
func (c *MessageController) AddReaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["id"]

	var reaction entity.Reaction
	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert string to uint
	messageIdUint, err := strconv.ParseUint(messageId, 10, 32)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
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
