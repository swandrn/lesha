package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lesha.com/server/internal/entity"
	"lesha.com/server/internal/services"
)

func main() {
	http.HandleFunc("/login", services.LoginHandler)
	http.HandleFunc("/protected", services.AuthMiddleware(services.ProtectedHandler))

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file ", err.Error())
	}
	dsn := os.Getenv("DB_URL")
	fmt.Println("dsn: ", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Migrating...")
	err = db.AutoMigrate(&entity.Channel{}, &entity.Friendship{}, &entity.Media{}, &entity.Message{}, &entity.Reaction{}, &entity.Server{}, &entity.User{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration successful!")

	
}
