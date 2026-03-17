package database

import (
	"database/sql"
	"time"
)

type User struct {
	ID             string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Email          string
	HashedPassword string
}

type Session struct {
	ID              string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	SessionID       string
	IsAuthenticated int32
	ExpiredAt       int64
	UserID          sql.NullString
}

type Menu struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Description string
	Price       float32
	IsStocked   int32
}
