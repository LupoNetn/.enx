package env

import (
	"net/http"

	"github.com/luponetn/enx/internal/db"
)

type authmiddleware func(handler http.HandlerFunc) http.Handler

func RegisterRoutes(mux *http.ServeMux, queries *db.Queries, protected authmiddleware) {
	service := NewService(queries)
	h := NewHandler(service)

	mux.Handle("POST /envs", protected(h.CreateEnv))
	mux.Handle("GET /envs/{id}", protected(h.GetEnvByID))
	mux.Handle("GET /projects/{project_id}/envs", protected(h.GetEnvsByProject))
	mux.Handle("PUT /envs/{id}", protected(h.UpdateEnv))
	mux.Handle("DELETE /envs/{id}", protected(h.DeleteEnv))
}
