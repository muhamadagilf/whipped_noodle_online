package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/muhamadagilf/whipped_noodle_online/util"
)

type transactionResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

func MidtransCreateTransaction(cart util.Cart, userPaymentDetail util.UserPaymentDetail) (transactionResponse, error) {
	var response transactionResponse
	type itemDetail struct {
		ID       string
		Name     string
		Quantity int
		Price    int
	}

	type customerDetail struct {
		FirstName, LastName string
		Email               string
		Phone               string
		BillingAddress      util.UserPaymentDetail
		ShippingAddress     util.UserPaymentDetail
	}

	type midtransPayload struct {
		ID             string
		GrossAmount    int32
		ItemsDetails   []itemDetail
		CustomerDetail customerDetail
	}

	var itemDetails []itemDetail
	for id, item := range cart.Menus {
		itemDetails = append(itemDetails, itemDetail{
			ID:       id,
			Name:     item.Name,
			Quantity: item.Qty,
			Price:    item.Price * item.Qty,
		})
	}

	payload := midtransPayload{
		ID:           cart.ID,
		GrossAmount:  cart.Total + cart.DeliveryFee,
		ItemsDetails: itemDetails,
		CustomerDetail: customerDetail{
			FirstName:       userPaymentDetail.Name,
			LastName:        userPaymentDetail.Name,
			Email:           userPaymentDetail.Email,
			Phone:           userPaymentDetail.Phone,
			BillingAddress:  userPaymentDetail,
			ShippingAddress: userPaymentDetail,
		},
	}

	JSONPayload, err := json.Marshal(payload)
	if err != nil {
		return response, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://app.sandbox.midtrans.com/snap/v1/transactions",
		bytes.NewReader(JSONPayload),
	)
	if err != nil {
		return response, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	authStr := base64.URLEncoding.EncodeToString([]byte(os.Getenv("PAYMENT_SERVER_KEY") + ":"))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", authStr))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return response, err
	}
	defer res.Body.Close()
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return response, err
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return response, err
	}

	return response, nil
}
