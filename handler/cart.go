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

func (h *Handler) AddToCartSession(c echo.Context) error {
	addToCartMenu := c.FormValue("menu-added")
	qtyStr := c.FormValue("menu-qty")
	s := strings.SplitN(addToCartMenu, ";", 3)
	menuID := s[0]
	menuName := s[1]
	priceStr := s[2]

	if menuName == "" || menuID == "" || priceStr == "" || qtyStr == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "empty request body")
	}
	if !slices.Contains(database.InMemoryMenu, menuName) {
		return echo.NewHTTPError(http.StatusBadRequest, util.NoMenuItem)
	}

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	qty, err := strconv.Atoi(qtyStr)
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
	if err = cart.Add(menuName, qty, price, menuID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	session.Values["cart"] = *cart
	if err = session.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.Response().Header().Set("HX-Redirect", "/home")
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) DeleteFromCartSession(c echo.Context) error {
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
	c.Response().Header().Set("HX-Redirect", "/home")
	return c.NoContent(http.StatusFound)
}
