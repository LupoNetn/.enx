package organization

import (
	"net/http"


	"github.com/luponetn/.enx/internal/db"
)

type authmiddleware  func(handler http.HandlerFunc) http.Handler


func RegisterRoutes(mux *http.ServeMux, queries *db.Queries, protected authmiddleware) {
     service := NewOrganizationService(queries)
	 h := NewOrganizationHandler(service)

	 mux.Handle("POST /organizations", protected(h.CreateOrganization))
	 mux.Handle("PUT /organizations/{id}", protected(h.UpdateOrganization))
	 mux.Handle("DELETE /organizations/{id}", protected(h.DeleteOrganization))
	 mux.Handle("GET /users/{user_id}/organizations", protected(h.GetAllOrganizationsByUser))
	 mux.Handle("GET /organizations/{id}/members", protected(h.GetAllUsersInOrganization))
}