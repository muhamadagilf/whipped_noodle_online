package middlewares

import "github.com/muhamadagilf/whipped_noodle_online/internal/server"

type Middlewares struct {
	Server *server.Server
}

func NewMiddlewares(server *server.Server) *Middlewares {
	return &Middlewares{Server: server}
}
