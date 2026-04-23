package test

import (
	"log"
	"testing"

	"github.com/muhamadagilf/whipped_noodle_online/service"
	"github.com/muhamadagilf/whipped_noodle_online/util"
	"github.com/stretchr/testify/assert"
)

func TestPaymentService(t *testing.T) {
	testCart := util.Cart{
		ID: "test-orderid-123",
		Menus: map[string]util.MenuOrder{
			"test-itemid-123": util.MenuOrder{
				Name:  "Test Nasi Goreng",
				Price: 5000,
				Qty:   2,
			},
		},
		TotalQty:    2,
		Total:       10000,
		DeliveryFee: 5000,
	}

	testPaymentDetail := util.UserPaymentDetail{
		Name:        "Testing Nama",
		Email:       "testingemail@example.com",
		Phone:       "08451123123",
		Address:     "Tunnel Winden no.33",
		City:        "Winden City",
		PostalCode:  "123123",
		CountryCode: "WDN",
	}

	responseTest, err := service.MidtransCreateTransaction(testCart, testPaymentDetail)
	assert.NoError(t, err, "no return error for midtrans result")
	assert.NotEmpty(t, responseTest, "no return empty midtrans result")
	assert.NotEmpty(t, responseTest.Token, "no return empty token")
	assert.NotEmpty(t, responseTest.RedirectURL, "no return redirect url")
	log.Println(responseTest)
}
