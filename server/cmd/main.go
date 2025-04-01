package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
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
