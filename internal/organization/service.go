package organization

import (
	"context"

	"github.com/luponetn/.enx/internal/db"
)

type Service struct {
	queries *db.Queries
}

type Svc interface {
	CreateOrganization(ctx context.Context, args db.CreateOrganizationParams) (db.CreateOrganizationRow, error)
}


func NewOrganizationService(queries *db.Queries) Svc {
  return &Service{
     queries: queries,
  }
}

//implement organization services
func (s *Service) CreateOrganization(ctx context.Context, args db.CreateOrganizationParams) (db.CreateOrganizationRow, error) {
	return s.queries.CreateOrganization(ctx, args)
}