package main

import (
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"time"
	"errors"
)

var (
	// Define a secret key for signing tokens
	jwtKey = []byte("your_secret_key")
)

// Token represents the structure of the authentication token
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// User represents a user in the system
type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Login represents the login request body
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Function to create a user (implement database logic here)
func createUser(email, password string) (string, error) {
	// Check if user exists and create new user logic here
	return "user_id", nil // return user ID or an error
}

// Function to generate tokens (implement logic here)
func tokensFromLogin(email, password string) (string, string, error) {
	// Validate user credentials and generate tokens
	accessToken := "generated_access_token"
	refreshToken := "generated_refresh_token"
	return accessToken, refreshToken, nil
}

// Middleware to check for token validity
func validateToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Check token signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// Set refresh cookie
func setRefreshCookie(c *gin.Context, refreshToken string) {
	c.SetCookie("refresh_token", refreshToken, 3600*24*30, "/", "", true, true) // 30 days
}

func main() {
	router := gin.Default()

	router.POST("/login", func(c *gin.Context) {
		var login Login
		if err := c.ShouldBindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		accessToken, refreshToken, err := tokensFromLogin(login.Email, login.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		data := Token{
			AccessToken:  accessToken,
			TokenType:    "bearer",
			RefreshToken: refreshToken,
		}
		setRefreshCookie(c, refreshToken)
		c.JSON(http.StatusOK, data)
	})

	router.POST("/register", func(c *gin.Context) {
		var login Login
		if err := c.ShouldBindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		userID, err := createUser(login.Email, login.Password)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}

		accessToken, refreshToken, err := tokensFromLogin(login.Email, login.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		data := Token{
			AccessToken:  accessToken,
			TokenType:    "bearer",
			RefreshToken: refreshToken,
		}
		setRefreshCookie(c, refreshToken)
		c.JSON(http.StatusOK, data)
	})

	router.GET("/refresh", func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "cookie not found"})
			return
		}

		// Here you would validate the refresh token and issue new tokens
		// For simplicity, just returning a dummy access token
		newAccessToken := "new_generated_access_token"
		c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
	})

	router.Run(":8080")
}
