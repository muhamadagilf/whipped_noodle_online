package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createUsersTable = `CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	verified_email TEXT NOT NULL UNIQUE,
	picture TEXT NOT NULL
);`

const createSessionsTable = `CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
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
	price INTEGER NOT NULL,
	is_stocked INTEGER NOT NULL DEFAULT 1
);`

const createOrdersTable = `CREATE TABLE IF NOT EXISTS orders (
	id TEXT PRIMARY KEY,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	menu_id TEXT REFERENCES menus(id) ON DELETE CASCADE
);`

const createCheckoutsTable = `CREATE TABLE IF NOT EXISTS checkouts (
	id TEXT PRIMARY KEY,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	order_id TEXT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	status TEXT NOT NULL,
	total_payment INTEGER NOT NULL
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

	if _, err := DB.Exec(createOrdersTable); err != nil {
		return err
	}

	if _, err := DB.Exec(createCheckoutsTable); err != nil {
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
			Name:        "mie seblak",
			Description: "Lorem Ipsum endurance testing desc type nooodle list",
			Price:       15000,
			IsStocked:   1,
		},
		{
			ID:          "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Name:        "mie goreng bumbu udang",
			Description: "Lorem Ipsum endurance testing desc type nooodle list",
			Price:       16000,
			IsStocked:   1,
		},
		{
			ID:          "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Name:        "mie rebus bumbu udang",
			Description: "Lorem Ipsum endurance testing desc type nooodle list",
			Price:       15000,
			IsStocked:   1,
		},
		{
			ID:          "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Name:        "nasi goreng",
			Description: "Lorem Ipsum endurance testing desc type nooodle list",
			Price:       16000,
			IsStocked:   1,
		},
	}
	if err := q.Transaction(context.Background(), myDB, func(qtx *Queries) error {
		for _, item := range items {
			id := uuid.New().String()
			_, err := myDB.Exec(insertMenus,
				id,
				item.Name,
				item.Description,
				item.Price,
				item.IsStocked,
			)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
