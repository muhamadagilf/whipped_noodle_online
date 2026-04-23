package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhamadagilf/whipped_noodle_online/middlewares"
	"github.com/muhamadagilf/whipped_noodle_online/util"
)

func (h *Handler) Checkoutpage(c echo.Context) error {
	query := h.Server.Queries
	csrf, ok := c.Get("csrf").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoCSRFError)
	}
	cart, ok := c.Get("cart").(*util.Cart)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoCartError)
	}
	cred, ok := c.Get("userCred").(middlewares.UserCred)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoUserIDError)
	}
	user, err := query.GetUserByEmail(c.Request().Context(), cred.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.Render(http.StatusOK, "checkout", Data{
		"csrf_token":   csrf,
		"cart":         cart,
		"user":         user,
		"totalPayment": cart.Total + cart.DeliveryFee,
	})
}
