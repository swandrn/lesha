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
	"lesha.com/server/internal/entity"
	"lesha.com/server/internal/repositories"
)

// JWT Claims structure
type Claims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid request format",
		})
		return
	}
	log.Printf("User: %+v", user)
	userRepository := repositories.NewUserRepository(database.Connect())
	existingUser, err := userRepository.GetUserByEmail(user.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Failed to check existing user",
		})
		return
	}
	if existingUser != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User with this email already exists",
		})
		return
	}
	err = userRepository.CreateUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Failed to create user",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
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
		allBlacklistedTokens, err := blacklistedTokenRepository.GetAllBlacklistedTokens()
		log.Printf("Token from request: '%s'", tokenString)
		log.Printf("Token from DB: '%s'", allBlacklistedTokens[0].Token)
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

	var creds entity.User
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userRepository := repositories.NewUserRepository(database.Connect())
	user, err := userRepository.GetUserByEmail(creds.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "No user found with this email",
		})
		return
	}

	if user.Password != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid credentials",
		})
		return
	}
	// Generate JWT
	token, err := GenerateJWT(fmt.Sprint(user.ID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Failed to generate token",
		})
		return
	}

	// Create a cookie containing the JWT
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		Path:     "/",
		MaxAge:   24 * 60 * 60, // 1 hour
		SameSite: http.SameSiteNoneMode,
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

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	db := database.Connect()
	blacklistedTokenRepository := repositories.NewBlacklistedTokenRepository(db)
	err = blacklistedTokenRepository.CreateBlacklistedToken(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to create blacklisted token", http.StatusInternalServerError)
		return
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUser", r)
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	tokenString := cookie.Value
	user, err := extractUserFromToken(tokenString)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User is connected",
		"user":    user,
	})

}

func extractUserFromToken(tokenString string) (*entity.User, error) {
	db := database.Connect()
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	userRepository := repositories.NewUserRepository(db)
	user, err := userRepository.GetUserById(claims.UserId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
