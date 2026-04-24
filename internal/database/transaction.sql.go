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

func (q *Queries) DeleteTransactionByID(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, "DELETE FROM transactions WHERE id=?;", id)
	return err
}

func (q *Queries) GetTransactionByID(ctx context.Context, id string) ([]Transaction, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT * FROM transactions WHERE id=?;", id)
	if err != nil {
		return nil, err
	}
	var s []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.MIdtransTransactionID,
			&i.Status,
			&i.TotalPayment,
		); err != nil {
			return nil, err
		}
		s = append(s, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return s, nil
}

func (q *Queries) GetTransactionByUserID(ctx context.Context, id string) ([]Transaction, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT * FROM transactions WHERE user_id=?;", id)
	if err != nil {
		return nil, err
	}
	var s []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Status,
			&i.TotalPayment,
			&i.UserID,
			&i.MIdtransTransactionID,
		); err != nil {
			return nil, err
		}
		s = append(s, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return s, nil
}
