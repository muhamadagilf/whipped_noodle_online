package handler

import (
	"errors"
	"slices"
	"strconv"
	"strings"

	"github.com/muhamadagilf/whipped_noodle_online/internal/database"
	"github.com/muhamadagilf/whipped_noodle_online/internal/server"
	"github.com/muhamadagilf/whipped_noodle_online/util"
)

type Handler struct {
	Server *server.Server
}

type (
	Data map[string]any
)

func NewHandler(server *server.Server) *Handler {
	return &Handler{Server: server}
}

type addToCart struct {
	menu string
	qty  int
}

func parseAddToCartRequest(m string) (addToCart, error) {
	splitM := strings.Split(m, ";")
	if len(splitM) != 2 {
		return addToCart{}, errors.New("invalid add_to_cart request")
	}
	qty, err := strconv.Atoi(splitM[1])
	if err != nil {
		return addToCart{}, err
	}
	splitM[0] = strings.ToLower(splitM[0])
	if !slices.Contains(database.InMemoryMenu, splitM[0]) {
		return addToCart{}, errors.New(util.NoMenuItem)
	}
	return addToCart{menu: splitM[0], qty: qty}, nil
}
