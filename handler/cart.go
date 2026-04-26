package handler

import (
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/muhamadagilf/whipped_noodle_online/internal/database"
	"github.com/muhamadagilf/whipped_noodle_online/util"
)

type cartMenu struct {
	MenuID   string `validate:"required"`
	MenuName string `validate:"required"`
	Price    string `validate:"required"`
	Qty      string `validate:"required"`
}

func (h *Handler) AddToCart(c echo.Context) error {
	addToCartMenu := c.FormValue("menu-added")
	s := strings.SplitN(addToCartMenu, ";", 3)
	menu := cartMenu{
		MenuID:   s[0],
		MenuName: s[1],
		Price:    s[2],
		Qty:      c.FormValue("menu-qty"),
	}

	if err := h.validate.Struct(menu); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if !slices.Contains(database.InMemoryMenu, menu.MenuName) {
		return echo.NewHTTPError(http.StatusBadRequest, util.NoMenuItem)
	}

	price, err := strconv.Atoi(menu.Price)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	qty, err := strconv.Atoi(menu.Qty)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	session, ok := c.Get("session").(*sessions.Session)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
	}
	cart, ok := c.Get("cart").(*util.Cart)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
	}
	if err = cart.Add(menu.MenuName, qty, int64(price), menu.MenuID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	session.Values["cart"] = *cart
	if err = session.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.Render(http.StatusOK, "add-alert", Data{})
	c.Render(http.StatusOK, "cart-count", Data{"cart": cart })
	return c.Render(http.StatusOK, "cart-menu-section", Data{ "cart": cart })
}

func (h *Handler) DeleteFromCart(c echo.Context) error {
	menuID := c.Param("menuID")
	session, ok := c.Get("session").(*sessions.Session)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
	}
	cart, ok := c.Get("cart").(*util.Cart)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
	}
	if err := cart.Remove(menuID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	session.Values["cart"] = *cart
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.Render(http.StatusOK, "cart-count", Data{"cart": cart})
	return c.Render(http.StatusOK, "cart-menu-section", Data{
		"cart": cart,
	})
}
