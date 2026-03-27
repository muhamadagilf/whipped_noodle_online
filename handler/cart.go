package handler

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/muhamadagilf/whipped_noodle_online/util"
)

func (h *Handler) AddToCartSession(c echo.Context) error {
	ctx := c.Request().Context()
	query := h.Server.Queries

	menuStr := c.FormValue("add-to-cart-menu")
	menu, err := parseAddToCartRequest(menuStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	session, ok := c.Get("session").(*sessions.Session)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
	}
	cart, ok := c.Get("cart").(*util.Cart)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
	}
	if err := cart.Add(ctx, query, menu.menu, menu.qty); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	session.Values["cart"] = *cart
	if err = session.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.Response().Header().Set("HX-Redirect", "/home")
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) DeleteFromCartSession(c echo.Context) error {
	ctx := c.Request().Context()
	query := h.Server.Queries
	menuStr := c.Param("menu")
	menuStr = strings.ToLower(menuStr)

	session, ok := c.Get("session").(*sessions.Session)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
	}
	cart, ok := c.Get("cart").(*util.Cart)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
	}
	if err := cart.Remove(ctx, query, menuStr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	session.Values["cart"] = *cart
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.Redirect(http.StatusFound, "/home")
}
