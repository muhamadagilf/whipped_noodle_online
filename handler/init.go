package handler

import (
	"github.com/muhamadagilf/whipped_noodle_online/internal/server"
)

type Handler struct {
	Server *server.Server
}

type Data map[string]any

func NewHandler(server *server.Server) *Handler {
	return &Handler{Server: server}
}
