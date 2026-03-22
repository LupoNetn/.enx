package main

import (
	"net/http"
	"time"

	"github.com/luponetn/.enx/internal/auth"
	"github.com/luponetn/.enx/internal/config"
	"github.com/luponetn/.enx/internal/db"
)

type Server struct {
	queries *db.Queries
	config  *config.Config
}

func NewServer(queries *db.Queries, cfg *config.Config) *Server {
	return &Server{queries: queries, config: cfg}
}

func (s *Server) Protected(handler http.HandlerFunc) http.Handler {
	return auth.AuthMiddleware(s.config.JWTAccessSecret, handler)
}

func CreateRouter(s *Server) *http.ServeMux {
	mux := http.NewServeMux()

	// register auth routes
	auth.RegisterRoutes(mux, s.queries, s.config)

	return mux
}

func StartServer(router *http.ServeMux, port string) error {
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		IdleTimeout:       time.Minute,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
	}
	return server.ListenAndServe()
}