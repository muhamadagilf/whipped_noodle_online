// Package server
package server

import (
	"database/sql"
	"errors"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/muhamadagilf/whipped_noodle_online/internal/database"
)

var PublicURL = []string{
	"/home",
	"/login",
	"/sign",
	"/auth/login",
	"/auth/oauth/callback",
	"/favicon.ico",
}

var ProtectedURL = []string{
	"/profile",
	"/checkout",
}

type Server struct {
	DB           *sql.DB
	Queries      *database.Queries
	SessionName  string
	SessionStore *sessions.CookieStore
}

func NewServer() (*Server, error) {
	db, err := sql.Open("sqlite", "./my.db")
	if err != nil {
		return nil, err
	}

	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		return nil, errors.New("cannot find SESSION_KEY in environment")
	}
	store := sessions.NewCookieStore([]byte(sessionKey))
	store.Options.Domain = ""
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false
	store.Options.SameSite = http.SameSiteStrictMode
	store.Options.MaxAge = 86400

	return &Server{
		DB:           db,
		Queries:      database.New(db),
		SessionName:  "web_session",
		SessionStore: store,
	}, nil
}
