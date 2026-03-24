package auth

import (
	"net/http"

	"github.com/luponetn/enx/internal/config"
	"github.com/luponetn/enx/internal/db"
)

func RegisterRoutes(mux *http.ServeMux, queries *db.Queries, cfg *config.Config) {
	service := NewAuthService(queries,cfg)
	h := NewAuthHandler(service)

	mux.HandleFunc("POST /auth/register", h.HandleRegister)
	mux.HandleFunc("POST /auth/login", h.HandleLogin)
	mux.HandleFunc("POST /auth/refresh", h.HandleRefreshToken)
}
