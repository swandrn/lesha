package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"lesha.com/server/internal/database"
	"lesha.com/server/internal/repositories"
)

// User structure (simulating a database)
type User struct {
	Id       string
	Username string
	Password string
}

// Dummy user data
var users = map[string]string{
	"admin": "password123",
}

// JWT Claims structure
type Claims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a JWT token for authenticated users
func GenerateJWT(userId string) (string, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file ", err.Error())
	}
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	return token.SignedString(jwtSecret)
}

// Middleware to validate JWT tokens
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from the Authorization header
		cookie, err := r.Cookie("token")

		if cookie == nil {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file ", err.Error())
		}
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))

		// Extract token from "Bearer <token>"
		tokenString := cookie.Value

		// Check if the token is blacklisted
		db := database.Connect()
		blacklistedTokenRepository := repositories.NewBlacklistedTokenRepository(db)
		blacklistedToken, err := blacklistedTokenRepository.GetBlacklistedToken(tokenString)
		if err != nil && err != gorm.ErrRecordNotFound {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		if blacklistedToken != nil {
			http.Error(w, "Token is blacklisted", http.StatusUnauthorized)
			return
		}

		// Parse the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		// Check if token is valid
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	}
}

// LoginHandler handles user authentication and returns a JWT
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var creds User
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Printf("creds: %+v\n", creds)

	if storedPassword, ok := users[creds.Username]; !ok || storedPassword != creds.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token, err := GenerateJWT(creds.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Create a cookie containing the JWT
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		Path:     "/",
		MaxAge:   3600, // 1 hour
	}

	// Set the cookie in the response
	http.SetCookie(w, cookie)

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully logged in",
	})
}

// ProtectedHandler is an example of a secured endpoint
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome! You have access to this protected route.")
}
