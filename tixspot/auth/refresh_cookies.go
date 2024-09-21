package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key")

// TokenClaims represents the structure for JWT claims
type TokenClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// newTokensUsingRefresh generates new access and refresh tokens using the refresh token
func newTokensUsingRefresh(refreshToken string) (map[string]string, error) {
	// Validate the refresh token (this should match your logic for refresh token validation)
	token, err := jwt.ParseWithClaims(refreshToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	// Generate new access and refresh tokens
	claims := token.Claims.(*TokenClaims)
	newAccessToken := generateToken(claims.Email, ACCESS_TOKEN_EXPIRE_MINUTES)
	newRefreshToken := generateToken(claims.Email, REFRESH_TOKEN_EXPIRE_HOURS)

	return map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}, nil
}

// generateToken creates a new JWT token with a specific expiration time
func generateToken(email string, expirationMinutes int64) string {
	expirationTime := time.Now().Add(time.Duration(expirationMinutes) * time.Minute)
	claims := &TokenClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}

// getRefreshCookie retrieves the refresh token from the request cookies
func getRefreshCookie(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Cookie not found"})
		return
	}
	tokens, err := newTokensUsingRefresh(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token"})
		return
	}
	c.JSON(http.StatusOK, tokens)
}

// setRefreshCookie sets the refresh token as a cookie
func setRefreshCookie(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	refreshToken, ok := data["refresh_token"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing refresh token"})
		return
	}

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(REFRESH_TOKEN_EXPIRE_HOURS)*3600,
		"/",
		"localhost",
		true,  // secure
		true,  // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{"message": "Cookie set successfully"})
}

func main() {
	router := gin.Default()

	router.POST("/get-refresh-cookie", getRefreshCookie)
	router.POST("/set-refresh-cookie", setRefreshCookie)

	router.Run(":8080")
}
