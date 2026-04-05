package util

import (
	"context"
	"errors"

	"github.com/muhamadagilf/whipped_noodle_online/internal/database"
)

var (
	NoSessionError     = "cannot find session in context"
	NoUserIDError      = "cannot find user_id in session"
	NoSessionIDError   = "cannot find session_id in session"
	NoCSRFError        = "cannot find csrf_token in context"
	NoCartError        = "cannot find cart in context"
	HTTPErrorAssertErr = "failed to assert type: echo.HTTPError"
	NoMenuItem         = "menu not found. please order in the menu"
)

type MenuOrder struct {
	Name       string
	Qty, Price int
}

type Cart struct {
	ID    string
	Menus map[string]MenuOrder
	Total int32
}

type Order struct {
	ID    string
	Menus []string
}

func (c *Cart) Add(menu string, qty int, price int, menuID string) error {
	if i, ok := c.Menus[menuID]; !ok {
		c.Menus[menuID] = MenuOrder{Name: menu, Qty: qty, Price: price}
	} else {
		i.Qty += qty
		c.Menus[menuID] = i
	}
	c.Total += int32(price * qty)
	return nil
}

func (c *Cart) Remove(menuID string) error {
	if i, ok := c.Menus[menuID]; ok {
		delete(c.Menus, menuID)
		c.Total -= int32(i.Price * i.Qty)
	} else {
		return errors.New("cannot find menu in the cart")
	}
	return nil
}

func (c *Cart) CreateOrder(
	ctx context.Context,
	query *database.Queries,
) error {
	return nil
}
