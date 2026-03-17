package database

import "context"

const getUserByEmail = `SELECT * FROM users WHERE email=?`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	if err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
	); err != nil {
		return i, err
	}
	if err := row.Err(); err != nil {
		return i, err
	}
	return i, nil
}

const createUser = `INSERT INTO users(id, name, email, verified_email, picture)
VALUES (?,?,?,?,?)`

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
		args.Name,
		args.Email,
		args.VerifiedEmail,
		args.Picture,
	)
	return err
}
