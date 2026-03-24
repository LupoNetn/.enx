package env

import (
	"net/http"

	"github.com/luponetn/enx/internal/db"
)

type authmiddleware  func(handler http.HandlerFunc) http.Handler

func RegisterRoutes(mux *http.ServeMux, queries *db.Queries, protected authmiddleware) {
	service := NewService(queries)
	h := NewHandler(service)

	mux.Handle("POST /env", protected(h.CreateEnv))
}
