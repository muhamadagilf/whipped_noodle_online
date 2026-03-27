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
	return c.Render(http.StatusOK, "checkout", Data{
		"csrf_token": csrf,
	})
}

// should new checkout entries created, early when user POST Request to /checkout
// or just creates new order entries first, and checkout later once the user hit "pay"
// the scenario very likely to lose the information easly, because it still might rely
// only on in-session information. to constructs orders information on the front-end
// however i would be lightweight operation for the system, because we dont rush to stores
// the order data to checkout table, cause user can just leaves it there in the cart
// wihtout even proceeds the payment.
// (or even not create entries on the order table, for more convinient DB optimization)
// main reason would be we dont have to bother to create another bg_worker to clean up orphange entries
// in both checkout and order table

func (h *Handler) Checkout(c echo.Context) error {
	return c.NoContent(http.StatusCreated)
}
