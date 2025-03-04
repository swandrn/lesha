package services

import (
	"fmt"

	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key for signing the JWT (keep it secure!)
var jwtSecret = []byte("your-secret-key")

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
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

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
	// Get username and password from query params (simplified for example)
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	// Validate user credentials
	if storedPassword, ok := users[username]; !ok || storedPassword != password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token, err := GenerateJWT(username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return token to the client
	fmt.Fprintf(w, "Token: %s", token)
}

// ProtectedHandler is an example of a secured endpoint
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome! You have access to this protected route.")
}

