package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"lesha.com/server/internal/entity"
	"lesha.com/server/internal/services"
)

type ChannelController struct {
	channelService *services.ChannelService
	serverService  *services.ServerService
}

func NewChannelController(channelService *services.ChannelService, serverService *services.ServerService) *ChannelController {
	return &ChannelController{
		channelService: channelService,
		serverService:  serverService,
	}
}

// GetChannels returns all channels
func (c *ChannelController) GetChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := c.channelService.GetChannels()
	if err != nil {
		http.Error(w, "Failed to fetch channels", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(channels)
}

// GetChannel returns a specific channel by ID
func (c *ChannelController) GetChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channelId := vars["id"]

	channel, err := c.channelService.GetChannel(channelId)
	if err != nil {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(channel)
}

// CreateChannel creates a new channel
func (c *ChannelController) CreateChannel(w http.ResponseWriter, r *http.Request) {
	var channel entity.Channel
	serverID, err := strconv.ParseUint(r.FormValue("serverID"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}
	channel.ServerID = uint(serverID)
	channel.Name = r.FormValue("name")

	// Log channel information for debugging
	fmt.Printf("Creating channel: Name=%s, ServerID=%d\n", channel.Name, channel.ServerID)

	if err := c.channelService.CreateChannel(&channel); err != nil {
		http.Error(w, "Failed to create channel", http.StatusInternalServerError)
		return
	}

	// Get all server members
	members, err := c.serverService.GetServerMembers(fmt.Sprintf("%d", serverID))
	if err != nil {
		http.Error(w, "Failed to get server members", http.StatusInternalServerError)
		return
	}

	// Add all members to the channel
	for _, member := range members {
		if err := c.channelService.AddUserToChannel(channel.ID, member.ID); err != nil {
			http.Error(w, "Failed to add user to channel", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(channel)
}

// UpdateChannel updates a channel's information
func (c *ChannelController) UpdateChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channelId := vars["id"]

	var updateData entity.Channel
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	channel, err := c.channelService.GetChannel(channelId)
	if err != nil {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	if updateData.Name != "" {
		channel.Name = updateData.Name
	}

	if err := c.channelService.UpdateChannel(channel); err != nil {
		http.Error(w, "Failed to update channel", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(channel)
}

// DeleteChannel deletes a channel
func (c *ChannelController) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channelId := vars["id"]

	channel, err := c.channelService.GetChannel(channelId)
	if err != nil {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
	}

	if err := c.channelService.DeleteChannel(channel); err != nil {
		http.Error(w, "Failed to delete channel", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Channel deleted successfully",
	})
}
