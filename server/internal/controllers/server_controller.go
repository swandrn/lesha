package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"lesha.com/server/internal/entity"
	"lesha.com/server/internal/services"
)

type ServerController struct {
	serverService *services.ServerService
}

func NewServerController(serverService *services.ServerService) *ServerController {
	return &ServerController{
		serverService: serverService,
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
	var server entity.Server
	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.serverService.CreateServer(&server); err != nil {
		http.Error(w, "Failed to create server", http.StatusInternalServerError)
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
