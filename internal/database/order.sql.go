package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createOrderQuery = `INSERT INTO orders (id, created_at, updated_at, qty, price, menu_id, transaction_id) VALUES (?,?,?,?,?,?,?);`

type CreateOrderParam struct {
	Qty           int
	Price         int64
	MenuID        string
	TransactionID string
}

func (q *Queries) CreateOrder(ctx context.Context, args CreateOrderParam) error {
	_, err := q.db.ExecContext(ctx, createOrderQuery,
		uuid.New(),
		time.Now().Local().String(),
		time.Now().Local().String(),
		args.Qty,
		args.Price,
		args.MenuID,
		args.TransactionID,
	)
	return err
}

func (q *Queries) GetOrdersByOrderID(ctx context.Context, orderID string) ([]Order, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT * FROM orders WHERE id=?;", orderID)
	if err != nil {
		return nil, err
	}
	var s []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.MenuID,
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

func (q *Queries) DeleteOrderByTransactionID(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, "DELETE FROM orders WHERE transaction_id=?;", id)
	return err
}
