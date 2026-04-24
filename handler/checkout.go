package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/muhamadagilf/whipped_noodle_online/internal/database"
	"github.com/muhamadagilf/whipped_noodle_online/middlewares"
	"github.com/muhamadagilf/whipped_noodle_online/util"
)

type checkoutHitory struct {
	Transaction database.Transaction
	Orders      []database.JoinOrderMenu
}

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

func (h *Handler) CheckoutHistory(c echo.Context) error {
	query := h.Server.Queries
	cred, ok := c.Get("userCred").(middlewares.UserCred)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
	}
	checkouts, err := query.GetTransactionByUserID(c.Request().Context(), cred.UserID.String)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var history []checkoutHitory
	for _, item := range checkouts {
		orders, err := query.GetJoinOrderByTransactionID(c.Request().Context(), item.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Join(errors.New("Error JoinOrderMenu"), err))
		}
		history = append(history, checkoutHitory{
			Transaction: item,
			Orders:      orders,
		})
	}

	return c.Render(http.StatusOK, "checkout-history", Data{
		"history": history,
	})
}
