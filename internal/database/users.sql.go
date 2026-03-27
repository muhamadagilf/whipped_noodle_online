package database

import (
	"context"
	"time"
)

const getUserByEmail = `SELECT * FROM users WHERE email=?`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	if err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Email,
		&i.VerifiedEmail,
		&i.Picture,
	); err != nil {
		return i, err
	}
	if err := row.Err(); err != nil {
		return i, err
	}
	return i, nil
}

const createUser = `INSERT INTO users(id, created_at, updated_at, name, email, verified_email, picture)
VALUES (?,?,?,?,?,?,?) ON CONFLICT (id) DO NOTHING;`

type CreateUserParam struct {
	ID            string
	Name          string
	Email         string
	VerifiedEmail bool
	Picture       string
}

func (q *Queries) CreateUser(ctx context.Context, args CreateUserParam) error {
	_, err := q.db.ExecContext(ctx, createUser,
		args.ID,
		time.Now().Local().String(),
		time.Now().Local().String(),
		args.Name,
		args.Email,
		args.VerifiedEmail,
		args.Picture,
	)
	return err
}
