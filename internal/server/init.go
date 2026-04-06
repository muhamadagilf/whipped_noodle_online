// Package server
package server

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/muhamadagilf/whipped_noodle_online/internal/database"
	"github.com/redis/go-redis/v9"
)

var PublicURL = []string{
	"/home",
	"/login",
	"/auth/login",
	"/auth/oauth/callback",
	"/favicon.ico",
	"/cart/add",
	"/cart/delete/:menuID",
}

var ProtectedURL = []string{
	"/auth/logout",
	"/profile",
	"/checkout",
}

type Server struct {
	DB           *sql.DB
	Queries      *database.Queries
	RDB          *redis.Client
	SessionName  string
	SessionStore *sessions.CookieStore
}

func NewServer() (*Server, error) {
	db, err := sql.Open("sqlite", "./my.db")
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		return nil, errors.New("cannot find SESSION_KEY in environment")
	}
	store := sessions.NewCookieStore([]byte(sessionKey))
	store.Options.Domain = ""
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false
	store.Options.SameSite = http.SameSiteLaxMode
	store.Options.MaxAge = int(12 * time.Hour)

	return &Server{
		DB:           db,
		Queries:      database.New(db),
		RDB:          rdb,
		SessionName:  "web_session",
		SessionStore: store,
	}, nil
}
