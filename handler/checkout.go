package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Checkout(c echo.Context) error {
	csrf, ok := c.Get("csrf").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot found csrf_token in context")
	}

	return c.Render(http.StatusOK, "checkout", Data{
		"csrf_token": csrf,
	})
}
