package database

import (
	"context"
	"database/sql"
	"time"
)

const createTransactionQuery = "INSERT INTO transactions (id, created_at, updated_at, status, total_payment, user_id, midtrans_transaction_id) VALUES (?,?,?,?,?,?,?);"

type CreateTransactionParam struct {
	ID           string
	UserID       string
	Status       string
	TotalPayment int64
}

func (q *Queries) CreateTransaction(ctx context.Context, args CreateTransactionParam) error {
	_, err := q.db.ExecContext(ctx, createTransactionQuery,
		args.ID,
		time.Now().Local(),
		time.Now().Local(),
		args.Status,
		args.TotalPayment,
		args.UserID,
		sql.NullString{Valid: false},
	)
	return err
}

const updateTransactionQuery = "UPDATE transactions SET status=?, midtrans_transaction_id=? WHERE id=?;"

type UpdateTransactionParam struct {
	Status, MID, ID string
}

func (q *Queries) UpdateTransactionByID(ctx context.Context, args UpdateTransactionParam) error {
	_, err := q.db.ExecContext(ctx, updateTransactionQuery, args.Status, args.MID, args.ID)
	return err
}
