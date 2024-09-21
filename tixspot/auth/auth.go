package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var (
	SECRET_KEY                   = []byte("09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7")
	ALGORITHM                    = "HS256"
	ACCESS_TOKEN_EXPIRE_MINUTES  = 30
	REFRESH_TOKEN_EXPIRE_HOURS    = 24 * 7
)

// Token structure for JWT
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// TokenData structure for token payload
type TokenData struct {
	Email string `json:"sub"`
}

// User structure
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Disabled bool   `json:"disabled"`
}

// UserInDB structure for user data with hashed password
type UserInDB struct {
	User
	HashedPassword string `json:"hashed_password"`
}

// HashPassword hashes the user's password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPassword compares a plain password with a hashed password
func VerifyPassword(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

// CreateAccessToken generates a new access token
func CreateAccessToken(email string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(ACCESS_TOKEN_EXPIRE_MINUTES) * time.Minute)
	claims := &TokenData{
		Email: email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// CreateRefreshToken generates a new refresh token
func CreateRefreshToken(email string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(REFRESH_TOKEN_EXPIRE_HOURS) * time.Hour)
	claims := &TokenData{
		Email: email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken verifies the token and returns the user email
func ValidateToken(tokenStr string, tokenType string) (string, error) {
	claims := &TokenData{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	if tokenType == "refresh" && claims.Email == "" {
		return "", errors.New("invalid refresh token")
	}
	return claims.Email, nil
}

// Main function
func main() {
	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		var user UserInDB
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		// Here you should verify the user with the database
		// For demonstration, we are assuming the user exists
		verifiedPassword := VerifyPassword("userPassword", user.HashedPassword) // Use hashed password
		if !verifiedPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect email or password"})
			return
		}

		accessToken, _ := CreateAccessToken(user.Email)
		refreshToken, _ := CreateRefreshToken(user.Email)

		c.JSON(http.StatusOK, Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "bearer",
		})
	})

	r.POST("/register", func(c *gin.Context) {
		var user UserInDB
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		// Hash the password
		hashedPassword, _ := HashPassword("password") // Replace with actual password
		user.HashedPassword = hashedPassword

		// Save user to the database (add your database logic here)

		c.JSON(http.StatusCreated, gin.H{"message": "user created"})
	})

	r.GET("/refresh", func(c *gin.Context) {
		refreshToken := c.Query("refresh_token")
		email, err := ValidateToken(refreshToken, "refresh")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}

		newAccessToken, _ := CreateAccessToken(email)
		newRefreshToken, _ := CreateRefreshToken(email)

		c.JSON(http.StatusOK, Token{
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
			TokenType:    "bearer",
		})
	})

	r.Run(":8080")
}
