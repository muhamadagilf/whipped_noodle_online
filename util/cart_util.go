package util

import "errors"

type MenuOrder struct {
	Name  string
	Price int64
	Qty   int
}

type Cart struct {
	ID          string
	Menus       map[string]MenuOrder
	TotalQty    int
	Total       int64
	DeliveryFee int64
}

type Order struct {
	ID    string
	Menus []string
}

func (c *Cart) Add(menu string, qty int, price int64, menuID string) error {
	if i, ok := c.Menus[menuID]; !ok {
		c.Menus[menuID] = MenuOrder{Name: menu, Qty: qty, Price: price}
	} else {
		i.Qty += qty
		c.Menus[menuID] = i
	}
	c.TotalQty += qty
	c.Total += price * int64(qty)
	return nil
}

func (c *Cart) Remove(menuID string) error {
	if i, ok := c.Menus[menuID]; ok {
		delete(c.Menus, menuID)
		c.Total -= i.Price * int64(i.Qty)
		c.TotalQty -= i.Qty
	} else {
		return errors.New("cannot find menu in the cart")
	}
	return nil
}
