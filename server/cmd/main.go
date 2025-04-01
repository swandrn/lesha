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
	// Auth routes
	r.HandleFunc("/login", services.LoginHandler).Methods("POST")
	r.HandleFunc("/register", services.RegisterHandler).Methods("POST")
	r.HandleFunc("/protected", services.AuthMiddleware(services.ProtectedHandler)).Methods("GET")
	r.HandleFunc("/logout", services.LogoutHandler).Methods("GET")
	r.HandleFunc("/get-user", services.GetUser).Methods("GET")
	r.HandleFunc("/ws", services.HandleWebSocket()).Methods("GET")

	userController := controllers.NewUserController(services.NewUserService(db))

	// User routes
	r.HandleFunc("/users", userController.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", userController.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", userController.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", userController.DeleteUser).Methods("DELETE")
	r.HandleFunc("/users/{id}/friends", userController.GetUserFriends).Methods("GET")

	// Initialize message controller
	messageController := controllers.NewMessageController(services.NewMessageService(db))

	// Message routes
	r.HandleFunc("/channels/{channelId}/messages", messageController.GetChannelMessages).Methods("GET")
	r.HandleFunc("/messages", messageController.CreateMessage).Methods("POST")
	r.HandleFunc("/messages/{id}", messageController.GetMessage).Methods("GET")
	r.HandleFunc("/messages/{id}/pin", messageController.PinMessage).Methods("POST")
	r.HandleFunc("/messages/{id}/reactions", messageController.AddReaction).Methods("POST")
	r.HandleFunc("/messages/{id}/reactions/{reactionId}", messageController.RemoveReaction).Methods("DELETE")

	// Initialize channel controller
	channelController := controllers.NewChannelController(services.NewChannelService(db))

	// Channel routes
	r.HandleFunc("/channels", channelController.GetChannels).Methods("GET")
	r.HandleFunc("/channels", channelController.CreateChannel).Methods("POST")
	r.HandleFunc("/channels/{id}", channelController.GetChannel).Methods("GET")
	r.HandleFunc("/channels/{id}", channelController.UpdateChannel).Methods("PUT")
	r.HandleFunc("/channels/{id}", channelController.DeleteChannel).Methods("DELETE")

	// Initialize server controller
	serverController := controllers.NewServerController(services.NewServerService(db), services.NewChannelService(db))

	// Server routes
	r.HandleFunc("/servers", serverController.GetServers).Methods("GET")

	r.HandleFunc("/servers", serverController.CreateServer).Methods("POST")
	r.HandleFunc("/servers/{id}", serverController.GetServer).Methods("GET")
	r.HandleFunc("/servers/{id}", serverController.UpdateServer).Methods("PUT")
	r.HandleFunc("/servers/{id}", serverController.DeleteServer).Methods("DELETE")

	// Setup CORS options
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // your frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(r)

	port := 8080
	fmt.Printf("Server running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))

}
