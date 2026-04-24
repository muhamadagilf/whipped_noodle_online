package test

import (
	"testing"

	"github.com/muhamadagilf/whipped_noodle_online/internal/database"
	"github.com/muhamadagilf/whipped_noodle_online/service"
	"github.com/muhamadagilf/whipped_noodle_online/util"
	"github.com/stretchr/testify/assert"
)

func TestPaymentService(t *testing.T) {
	testCart := util.Cart{
		ID: "test-orderid-456",
		Menus: map[string]util.MenuOrder{
			"test-itemid-456": util.MenuOrder{
				Name:  "Test Nasi Goreng",
				Price: 5000,
				Qty:   2,
			},
		},
		TotalQty:    2,
		Total:       10000,
		DeliveryFee: 5000,
	}

	testPaymentDetail := database.UserPaymentDetail{
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
}
