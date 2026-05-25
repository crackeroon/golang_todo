package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/crackeroon/golang_todo/internal/core/transport/http/middleware"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []core_http_middleware.MiddleWare
}

func (r *Route) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleWare(
		r.Handler,
		r.Middleware...,
	)
}
