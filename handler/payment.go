package handler

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/muhamadagilf/whipped_noodle_online/service"
	"github.com/muhamadagilf/whipped_noodle_online/util"
)

func (h *Handler) Pay(c echo.Context) error {
	formPaymentData := util.UserPaymentDetail{
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

	response, err := service.MidtransCreateTransaction(*cart, formPaymentData)
	if err != nil {
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

	// PROCESS

	return c.NoContent(http.StatusOK)
}

func (h *Handler) TransactionSuccess(c echo.Context) error {
	return c.Render(http.StatusOK, "payment-success", Data{})
}

func (h *Handler) TransactionFailed(c echo.Context) error {
	return c.Render(http.StatusOK, "payment-failed", Data{})
}
