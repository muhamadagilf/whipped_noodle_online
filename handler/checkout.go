package handler

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/muhamadagilf/whipped_noodle_online/util"
)

func (h *Handler) Checkoutpage(c echo.Context) error {
	csrf, ok := c.Get("csrf").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoCSRFError)
	}
	return c.Render(http.StatusOK, "checkout", Data{
		"csrf_token": csrf,
	})
}

// so, go with in-session checkout items information. and proceeded with DB operation, once the user hit the payment
// to really stored the data in DB
func (h *Handler) Checkout(c echo.Context) error {
	session, ok := c.Get("session").(*sessions.Session)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
	}
	// userCred, ok := c.Get("userCred").(middlewares.UserCred)
	// if !ok {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, util.NoUserIDError)
	// }
	// cart, ok := c.Get("cart").(*util.Cart)
	// if !ok {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, util.NoCartError)
	// }

	if err := session.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	c.Response().Header().Set("HX-Redirect", "/checkout")
	return c.NoContent(http.StatusCreated)
}
