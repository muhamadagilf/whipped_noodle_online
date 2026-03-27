package database

import (
	"database/sql"
	"time"
)

type User struct {
	ID            string
	CreatedAt     string
	UpdatedAt     string
	Name          string
	Email         string
	VerifiedEmail string
	Picture       string
}

type Session struct {
	ID              string
	CreatedAt       string
	UpdatedAt       string
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
	Price       int
	IsStocked   int32
}

type Order struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	MenuID    string
}

type Checkout struct {
	ID           string
	CreatedAt    string
	UpdatedAt    string
	OrderID      string
	UserID       string
	Status       string
	TotalPayment int32
}
