package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/muhamadagilf/whipped_noodle_online/internal/server"
)

type Handler struct {
	Server   *server.Server
	validate *validator.Validate
}

type (
	Data map[string]any
)

func NewHandler(server *server.Server, validate *validator.Validate) *Handler {
	return &Handler{Server: server, validate: validate}
}

type addToCart struct {
	menu string
	qty  int
}
