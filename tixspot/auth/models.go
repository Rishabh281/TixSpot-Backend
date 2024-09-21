package main

// Token represents the structure of the authentication token
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
