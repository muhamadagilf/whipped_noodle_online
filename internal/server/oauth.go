package server

import (
	"errors"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig *oauth2.Config

func InitOAuth() error {
	clientID := os.Getenv("OAUTH_CLIENT_ID")
	oauthClientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
	if clientID == "" || oauthClientSecret == "" {
		return errors.New("oauth client_id or secret not found")
	}
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: oauthClientSecret,
		RedirectURL:  "http://localhost:8000/auth/oauth/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return nil
}
