package service

import (
	"errors"
	"log"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/muhamadagilf/whipped_noodle_online/util"
)

type transactionResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

func MidtransCreateTransaction(cart util.Cart, userPaymentDetail util.UserPaymentDetail) (*snap.Response, error) {
	var itemDetails []midtrans.ItemDetails
	for id, item := range cart.Menus {
		itemDetails = append(itemDetails, midtrans.ItemDetails{
			ID:    id,
			Name:  item.Name,
			Price: int64(item.Price * item.Qty),
			Qty:   int32(item.Qty),
		})
		log.Println("[TRANSACTION DEBUG] ItemDetail", int64(item.Price*item.Qty))
	}

	serverKey := os.Getenv("PAYMENT_SERVER_KEY")
	if serverKey == "" {
		return nil, errors.New("invalid API Server Key. empty key")
	}
	s := snap.Client{}
	s.New(serverKey, midtrans.Sandbox)
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  cart.ID,
			GrossAmt: int64(cart.Total),
		},
		Items: &itemDetails,
		CustomerDetail: &midtrans.CustomerDetails{
			FName: userPaymentDetail.Name,
			LName: userPaymentDetail.Name,
			Email: userPaymentDetail.Email,
			Phone: userPaymentDetail.Phone,
			ShipAddr: &midtrans.CustomerAddress{
				FName:       userPaymentDetail.Name,
				LName:       userPaymentDetail.Name,
				Phone:       userPaymentDetail.Phone,
				Address:     userPaymentDetail.Address,
				City:        userPaymentDetail.City,
				Postcode:    userPaymentDetail.PostalCode,
				CountryCode: userPaymentDetail.CountryCode,
			},
		},
	}

	log.Println("[TRANSACTION DEBUG] GrossAmount:", int64(cart.Total))

	response, err := s.CreateTransaction(req)
	if err != nil {
		return response, err
	}

	return response, nil
}
