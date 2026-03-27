package database

import (
	"context"
	"time"
)

const createOrderQuery = `INSERT INTO orders (id, created_at, updated_at, menu_id) VALUES (?,?,?,?);`

func (q *Queries) CreateOrder(ctx context.Context, orderID string, menuID string) error {
	_, err := q.db.ExecContext(ctx, createOrderQuery,
		orderID,
		time.Now().Local().String(),
		time.Now().Local().String(),
		menuID,
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
