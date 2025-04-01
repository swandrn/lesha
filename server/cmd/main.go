package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"lesha.com/server/internal/controllers"
	"lesha.com/server/internal/database"
	"lesha.com/server/internal/entity"
	"lesha.com/server/internal/services"
	"lesha.com/server/internal/ws"
)

func main() {

	//Migrations
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file ", err.Error())
	}
	db := database.Connect()
	fmt.Println("Migrating...")
	err = db.AutoMigrate(&entity.Channel{}, &entity.Friendship{}, &entity.Media{}, &entity.Message{}, &entity.Reaction{}, &entity.Server{}, &entity.User{}, &entity.BlacklistedToken{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration successful!")

	r := mux.NewRouter()

	// Serve static files from uploads directory
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Auth routes
	r.HandleFunc("/login", services.LoginHandler).Methods("POST")
	r.HandleFunc("/register", services.RegisterHandler).Methods("POST")
	r.HandleFunc("/protected", services.AuthMiddleware(services.ProtectedHandler)).Methods("GET")
	r.HandleFunc("/logout", services.LogoutHandler).Methods("GET")
	r.HandleFunc("/get-user", services.GetUser).Methods("GET")
	r.HandleFunc("/ws", ws.HandleWebSocket(db)).Methods("GET")

	userController := controllers.NewUserController(services.NewUserService(db))

	// User routes
	r.HandleFunc("/users", services.AuthMiddleware(userController.GetUsers)).Methods("GET")
	r.HandleFunc("/users/{id}", services.AuthMiddleware(userController.GetUser)).Methods("GET")
	r.HandleFunc("/users/{id}", services.AuthMiddleware(userController.UpdateUser)).Methods("PUT")
	r.HandleFunc("/users/{id}", services.AuthMiddleware(userController.DeleteUser)).Methods("DELETE")
	r.HandleFunc("/users/{id}/friends", services.AuthMiddleware(userController.GetUserFriends)).Methods("GET")

	// Initialize message controller
	messageController := controllers.NewMessageController(services.NewMessageService(db))

	// Message routes
	r.HandleFunc("/channels/{channelID}/messages", services.AuthMiddleware(messageController.GetChannelMessages)).Methods("GET")
	r.HandleFunc("/messages", services.AuthMiddleware(messageController.CreateMessage)).Methods("POST")
	r.HandleFunc("/messages/{id}", services.AuthMiddleware(messageController.GetMessage)).Methods("GET")
	r.HandleFunc("/messages/{id}/pin", services.AuthMiddleware(messageController.PinMessage)).Methods("GET")
	r.HandleFunc("/messages/{id}/unpin", services.AuthMiddleware(messageController.UnpinMessage)).Methods("GET")
	r.HandleFunc("/messages/{id}/reactions", services.AuthMiddleware(messageController.AddReaction)).Methods("POST")
	r.HandleFunc("/messages/{id}/reactions/{reactionId}", services.AuthMiddleware(messageController.RemoveReaction)).Methods("DELETE")
	r.HandleFunc("/messages/{id}/media", services.AuthMiddleware(messageController.AddMedia)).Methods("POST")

	// Initialize channel controller
	channelController := controllers.NewChannelController(services.NewChannelService(db), services.NewServerService(db))

	// Channel routes
	r.HandleFunc("/channels", services.AuthMiddleware(channelController.GetChannels)).Methods("GET")
	r.HandleFunc("/channels", services.AuthMiddleware(channelController.CreateChannel)).Methods("POST")
	r.HandleFunc("/channels/{id}", services.AuthMiddleware(channelController.GetChannel)).Methods("GET")
	r.HandleFunc("/channels/{id}", services.AuthMiddleware(channelController.UpdateChannel)).Methods("PUT")
	r.HandleFunc("/channels/{id}", services.AuthMiddleware(channelController.DeleteChannel)).Methods("DELETE")
	r.HandleFunc("/servers/{id}/channels", services.AuthMiddleware(channelController.GetServerChannels)).Methods("GET")

	// Initialize server controller
	serverController := controllers.NewServerController(services.NewServerService(db), services.NewChannelService(db))

	// Server routes
	r.HandleFunc("/servers", services.AuthMiddleware(serverController.GetUserServers)).Methods("GET")
	r.HandleFunc("/servers", services.AuthMiddleware(serverController.CreateServer)).Methods("POST")
	r.HandleFunc("/servers/{id}", services.AuthMiddleware(serverController.GetServer)).Methods("GET")
	r.HandleFunc("/servers/{id}", services.AuthMiddleware(serverController.UpdateServer)).Methods("PUT")
	r.HandleFunc("/servers/{id}", services.AuthMiddleware(serverController.DeleteServer)).Methods("DELETE")
	r.HandleFunc("/servers/{id}/add-user", services.AuthMiddleware(serverController.AddUserToServerByEmail)).Methods("POST")

	// Setup CORS options
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5174"}, // your frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(r)

	port := 8080
	fmt.Printf("Server running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))

}
