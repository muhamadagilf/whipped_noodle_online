package database

import "context"

var InMemoryMenu = []string{
	"mie seblak",
	"nasi goreng",
	"mie goreng bumbu udang",
	"mie rebus bumbu udang",
}

func (q *Queries) GetAllMenu(ctx context.Context) ([]Menu, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT * FROM menus;")
	if err != nil {
		return nil, err
	}
	var menu []Menu
	for rows.Next() {
		var i Menu
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.IsStocked,
		); err != nil {
			return nil, err
		}
		menu = append(menu, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return menu, nil
}

func (q *Queries) GetMenuByName(ctx context.Context, name string) (Menu, error) {
	row := q.db.QueryRowContext(ctx, "SELECT * FROM menus WHERE name=?", name)
	var i Menu
	if err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.IsStocked,
	); err != nil {
		return i, err
	}
	if err := row.Err(); err != nil {
		return i, err
	}
	return i, nil
}
