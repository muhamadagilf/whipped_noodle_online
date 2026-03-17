package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

const createUsersTable = `CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	verified_email TEXT NOT NULL UNIQUE,
	picture TEXT NOT NULL
);`

const createSessionsTable = `CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	session_id TEXT NOT NULL UNIQUE,
	is_authenticated INTEGER NOT NULL DEFAULT 0,
	expired_at INTEGER NOT NULL,
	user_id TEXT REFERENCES users(id) ON DELETE CASCADE
);`

const createItemsTable = `CREATE TABLE IF NOT EXISTS menus (
	id TEXT PRIMARY KEY,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	name TEXT NOT NULL UNIQUE,
	description TEXT NOT NULL,
	price FLOAT NOT NULL,
	is_stocked INTEGER NOT NULL DEFAULT 1
);`

func Migrate(DB *sql.DB) error {
	if _, err := DB.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return err
	}

	if _, err := DB.Exec(createUsersTable); err != nil {
		return err
	}

	if _, err := DB.Exec(createSessionsTable); err != nil {
		return err
	}

	if _, err := DB.Exec(createItemsTable); err != nil {
		return err
	}

	return nil
}

func (q *Queries) InsertMainMenuToDB() error {
	myDB := q.db.(*sql.DB)
	insertMenus := "INSERT INTO menus (id, name, description, price, is_stocked) VALUES (?, ?, ?, ?, ?) ON CONFLICT (name) DO NOTHING;"
	items := []Menu{
		{
			ID:          "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Name:        "Mie Seblak",
			Description: "Lorem Ipsum endurance testing desc type nooodle list",
			Price:       15.000,
			IsStocked:   1,
		},
		{
			ID:          "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Name:        "Mie Goreng Bumbu Udang",
			Description: "Lorem Ipsum endurance testing desc type nooodle list",
			Price:       16.000,
			IsStocked:   1,
		},
		{
			ID:          "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Name:        "Mie Rebus Bumbu Udang",
			Description: "Lorem Ipsum endurance testing desc type nooodle list",
			Price:       15.000,
			IsStocked:   1,
		},
		{
			ID:          "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Name:        "Nasi Goreng",
			Description: "Lorem Ipsum endurance testing desc type nooodle list",
			Price:       16.000,
			IsStocked:   1,
		},
	}
	if err := q.Transaction(context.Background(), myDB, func(qtx *Queries) error {
		for _, item := range items {
			id := uuid.New().String()
			res, err := myDB.Exec(insertMenus,
				id,
				item.Name,
				item.Description,
				item.Price,
				item.IsStocked,
			)
			if err != nil {
				return err
			}
			log.Println("[INFO_DB]: ", res)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
