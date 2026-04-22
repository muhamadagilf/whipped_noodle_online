package util

import "errors"

type UserPaymentDetail struct {
	Name        string `validate:"required"`
	Email       string `validate:"required"`
	Phone       string `validate:"required"`
	Address     string `validate:"required"`
	City        string `validate:"required"`
	PostalCode  string `validate:"required"`
	CountryCode string `validate:"required"`
}

type MenuOrder struct {
	Name       string
	Qty, Price int
}

type Cart struct {
	ID          string
	Menus       map[string]MenuOrder
	TotalQty    int
	Total       int32
	DeliveryFee int32
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
	c.TotalQty += qty
	c.Total += int32(price * qty)
	return nil
}

func (c *Cart) Remove(menuID string) error {
	if i, ok := c.Menus[menuID]; ok {
		delete(c.Menus, menuID)
		c.Total -= int32(i.Price * i.Qty)
		c.TotalQty -= i.Qty
	} else {
		return errors.New("cannot find menu in the cart")
	}
	return nil
}
