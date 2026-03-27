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

type Cart struct {
	ID    string
	Menus map[string]int
	Total int32
}

func (c *Cart) Add(
	ctx context.Context,
	query *database.Queries,
	menu string, qty int,
) error {
	if _, ok := c.Menus[menu]; !ok {
		c.Menus[menu] = qty
	} else {
		c.Menus[menu] += qty
	}

	menuData, err := query.GetMenuByName(ctx, menu)
	if err != nil {
		return err
	}
	c.Total += int32(menuData.Price * qty)
	return nil
}

func (c *Cart) Remove(
	ctx context.Context,
	query *database.Queries,
	menu string,
) error {
	if qty, ok := c.Menus[menu]; ok {
		delete(c.Menus, menu)
		menuData, err := query.GetMenuByName(ctx, menu)
		if err != nil {
			return err
		}
		c.Total -= int32(menuData.Price * qty)
	} else {
		return errors.New("cannot find menu in the cart")
	}
	return nil
}
