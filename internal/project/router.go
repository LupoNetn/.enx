package project

import (
	"net/http"

	"github.com/luponetn/enx/internal/db"
)

type authmiddleware func(handler http.HandlerFunc) http.Handler

func RegisterRoutes(mux *http.ServeMux, queries *db.Queries, protected authmiddleware) {
	service := NewProjectService(queries)
	h := NewProjectHandler(service)

	mux.Handle("POST /projects", protected(h.CreateProject))
	mux.Handle("PUT /projects/{id}", protected(h.UpdateProject))
	mux.Handle("DELETE /projects/{id}", protected(h.DeleteProject))
	mux.Handle("GET /users/{user_id}/projects", protected(h.GetProjectsByUser))
	mux.Handle("GET /projects/{id}/members", protected(h.GetAllUsersInProject))
	mux.Handle("GET /projects/name/{name}", protected(h.GetProjectByName))
}
