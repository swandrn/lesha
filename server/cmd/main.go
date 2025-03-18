package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"lesha.com/server/internal/database"
	"lesha.com/server/internal/entity"
	"lesha.com/server/internal/services"
)

func main() {

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

	// Start server

	r.HandleFunc("/login", services.LoginHandler).Methods("POST")
	r.HandleFunc("/protected", services.AuthMiddleware(services.ProtectedHandler)).Methods("GET")
	r.HandleFunc("/logout", services.LogoutHandler).Methods("POST")

	port := 8080
	fmt.Printf("Server running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))

}
