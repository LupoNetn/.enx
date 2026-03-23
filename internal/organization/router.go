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
}