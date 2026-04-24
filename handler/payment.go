package handler

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/muhamadagilf/whipped_noodle_online/internal/database"
	"github.com/muhamadagilf/whipped_noodle_online/middlewares"
	"github.com/muhamadagilf/whipped_noodle_online/service"
	"github.com/muhamadagilf/whipped_noodle_online/util"
)

func (h *Handler) Pay(c echo.Context) error {
	query := h.Server.Queries
	formPaymentData := database.UserPaymentDetail{
		Name:        c.FormValue("fullname"),
		Email:       c.FormValue("email"),
		Phone:       c.FormValue("phone"),
		Address:     c.FormValue("address"),
		City:        "Kota Cirebon",
		PostalCode:  c.FormValue("postal_code"),
		CountryCode: "IDN",
	}

	if err := h.validate.Struct(formPaymentData); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	cart, ok := c.Get("cart").(*util.Cart)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoCartError)
	}

	cred, ok := c.Get("userCred").(middlewares.UserCred)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoUserIDError)
	}

	err := query.Transaction(c.Request().Context(), h.Server.DB, func(qtx *database.Queries) error {
		var totalPayment int64
		transactionID := fmt.Sprintf("TRA-%v", uuid.New())
		if err := qtx.CreateTransaction(c.Request().Context(), database.CreateTransactionParam{
			ID:           transactionID,
			UserID:       cred.UserID.String,
			Status:       "PENDING",
			TotalPayment: cart.Total + cart.DeliveryFee,
		}); err != nil {
			return err
		}

		for id, item := range cart.Menus {
			menu, err := qtx.GetMenuByID(c.Request().Context(), id)
			if err != nil {
				return err
			}

			totalPayment += menu.Price * int64(item.Qty)

			if err := qtx.CreateOrder(c.Request().Context(), database.CreateOrderParam{
				Qty:           item.Qty,
				Price:         menu.Price,
				MenuID:        menu.ID,
				TransactionID: transactionID,
			}); err != nil {
				return err
			}
		}

		log.Println("REACH HERE AFTER INSERT ORDERS")
		if totalPayment != cart.Total {
			if err := qtx.DeleteTransactionByID(c.Request().Context(), transactionID); err != nil {
				return err
			}
			if err := qtx.DeleteOrderByTransactionID(c.Request().Context(), transactionID); err != nil {
				return err
			}
			return util.InvalidTotalPayment
		}

		cart.ID = transactionID

		return nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	log.Println("[CART_DEBUG]", cart.ID)

	response, err := service.MidtransCreateTransaction(*cart, formPaymentData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	session, ok := c.Get("session").(*sessions.Session)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, util.NoSessionError)
	}
	delete(session.Values, "cart")
	if err := session.Save(c.Request(), c.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.Response().Header().Set("HX-Redirect", response.RedirectURL)
	return c.NoContent(http.StatusFound)

}

func (h *Handler) TransactionNotification(c echo.Context) error {
	type notificationRequestBody struct {
		TransactionType          string `json:"transaction_type"`
		TransactionTime          string `json:"transaction_time"`
		TransactionStatus        string `json:"transaction_status"`
		TransactionID            string `json:"transaction_id"`
		StatusMessage            string `json:"status_message"`
		StatusCode               string `json:"status_code"`
		SignatureKey             string `json:"signature_key"`
		SettlementTime           string `json:"settlement_time"`
		PaymentType              string `json:"payment_type"`
		OrderID                  string `json:"order_id"`
		MerchantID               string `json:"merchant_id"`
		MerchantCrossReferenceID string `json:"merchant_cross_reference_id"`
		Issuer                   string `json:"issuer"`
		GrossAmount              string `json:"gross_amount"`
		FraudStatus              string `json:"fraud_status"`
		Currency                 string `json:"currency"`
		Acquirer                 string `json:"acquirer"`
	}

	payload := notificationRequestBody{}
	request := c.Request()
	defer request.Body.Close()

	reqData, err := io.ReadAll(request.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if err := json.Unmarshal(reqData, &payload); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	notifSigStr := payload.OrderID + payload.StatusCode + payload.GrossAmount + os.Getenv("PAYMENT_SERVER_KEY")
	hash := sha512.New()
	if _, err = hash.Write([]byte(notifSigStr)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	notifSignature := hex.EncodeToString(hash.Sum(nil))
	if notifSignature != payload.SignatureKey {
		return echo.NewHTTPError(http.StatusBadRequest, util.InvalidNotificationSignatureKey)
	}

	query := h.Server.Queries
	switch payload.TransactionStatus {
	case "settlement", "capture":
		if err := query.UpdateTransactionByID(c.Request().Context(), database.UpdateTransactionParam{
			Status: "PAID",
			MID:    payload.TransactionID,
			ID:     payload.OrderID,
		}); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	case "pending":
		if err := query.UpdateTransactionByID(c.Request().Context(), database.UpdateTransactionParam{
			Status: "PENDING",
			MID:    payload.TransactionID,
			ID:     payload.OrderID,
		}); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	case "expire", "cancel", "deny":
		if err := query.UpdateTransactionByID(c.Request().Context(), database.UpdateTransactionParam{
			Status: "FAILED",
			MID:    payload.TransactionID,
			ID:     payload.OrderID,
		}); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) TransactionSuccess(c echo.Context) error {
	return c.Render(http.StatusOK, "payment-success", Data{})
}

func (h *Handler) TransactionFailed(c echo.Context) error {
	return c.Render(http.StatusOK, "payment-failed", Data{})
}
