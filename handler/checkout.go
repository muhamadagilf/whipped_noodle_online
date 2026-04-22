package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhamadagilf/whipped_noodle_online/util"
)

func (h *Handler) Checkoutpage(c echo.Context) error {
	csrf, ok := c.Get("csrf").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoCSRFError)
	}
	cart, ok := c.Get("cart").(*util.Cart)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoCartError)
	}
	return c.Render(http.StatusOK, "checkout", Data{
		"csrf_token":   csrf,
		"cart":         cart,
		"totalPayment": cart.Total + cart.DeliveryFee,
	})
}
