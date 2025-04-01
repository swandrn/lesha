package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
	"lesha.com/server/internal/entity"
	"lesha.com/server/internal/services"
)

type ServerController struct {
	serverService  *services.ServerService
	channelService *services.ChannelService
}

func NewServerController(serverService *services.ServerService, channelService *services.ChannelService) *ServerController {
	return &ServerController{
		serverService:  serverService,
		channelService: channelService,
	}
}

// GetServers returns all servers
func (c *ServerController) GetServers(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uint)
	servers, err := c.serverService.GetUserServers(userID)
	if err != nil {
		http.Error(w, "Failed to fetch servers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(servers)
}

// GetServer returns a specific server by ID
func (c *ServerController) GetServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serverId := vars["id"]

	server, err := c.serverService.GetServer(serverId)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(server)
}

// CreateServer creates a new server
func (c *ServerController) CreateServer(w http.ResponseWriter, r *http.Request) {
	// Log the request
	fmt.Printf("Creating server request from %s: %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)

	// Extract userID from JWT token in context
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

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	// Get the file from form
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create uploads directory if it doesn't exist
	uploadDir := "uploads/servers"
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

	// Create server with image path
	server := entity.Server{
		Name:        name,
		Description: description,
		Image:       filepath,
		UserID:      userID,
	}

	if err := c.serverService.CreateServer(&server); err != nil {
		http.Error(w, "Failed to create server", http.StatusInternalServerError)
		return
	}

	// Add the creator as a member of the server
	if err := c.serverService.AddUserToServer(server.ID, userID); err != nil {
		http.Error(w, "Failed to add user to server", http.StatusInternalServerError)
		return
	}

	// Create default text channel
	channel := entity.Channel{
		Name:     "text",
		ServerID: server.ID,
	}
	if err := c.channelService.CreateChannel(&channel); err != nil {
		http.Error(w, "Failed to create default channel", http.StatusInternalServerError)
		return
	}

	// Add the creator to the channel
	if err := c.channelService.AddUserToChannel(channel.ID, userID); err != nil {
		http.Error(w, "Failed to add user to channel", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(server)
}

// UpdateServer updates a server's information
func (c *ServerController) UpdateServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serverId := vars["id"]

	var updateData entity.Server
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	server, err := c.serverService.GetServer(serverId)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}

	// Update fields
	if updateData.Name != "" {
		server.Name = updateData.Name
	}
	if updateData.Description != "" {
		server.Description = updateData.Description
	}
	if updateData.Image != "" {
		server.Image = updateData.Image
	}

	if err := c.serverService.UpdateServer(server); err != nil {
		http.Error(w, "Failed to update server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(server)
}

// DeleteServer deletes a server
func (c *ServerController) DeleteServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serverId := vars["id"]

	server, err := c.serverService.GetServer(serverId)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}

	if err := c.serverService.DeleteServer(server); err != nil {
		http.Error(w, "Failed to delete server", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Server deleted successfully",
	})
}

// AddUserToServerByEmail adds a user to a server using their email
func (c *ServerController) AddUserToServerByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serverId := vars["id"]

	var data struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the server first
	server, err := c.serverService.GetServer(serverId)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}

	// First get the user by email
	user, err := c.serverService.GetUserByEmail(data.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Add user to server
	if err := c.serverService.AddUserToServer(server.ID, user.ID); err != nil {
		http.Error(w, "Failed to add user to server", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User added to server successfully",
	})
}
