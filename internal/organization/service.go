package organization

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/.enx/internal/db"
)

type Service struct {
	queries *db.Queries
}

type Svc interface {
	CreateOrganization(ctx context.Context, args db.CreateOrganizationParams) (db.CreateOrganizationRow, error)
	UpdateOrganization(ctx context.Context, args db.UpdateOrganizationParams) (db.UpdateOrganizationRow, error)
	DeleteOrganization(ctx context.Context, id pgtype.UUID) error
	GetAllOrganizationsByUser(ctx context.Context, userID pgtype.UUID) ([]db.GetAllOrganizationsByUserRow, error)
	GetAllUsersInOrganization(ctx context.Context, organizationID pgtype.UUID) ([]db.GetAllUsersInOrganizationRow, error)
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

func (s *Service) UpdateOrganization(ctx context.Context, args db.UpdateOrganizationParams) (db.UpdateOrganizationRow, error) {
	return s.queries.UpdateOrganization(ctx, args)
}

func (s *Service) DeleteOrganization(ctx context.Context, id pgtype.UUID) error {
	return s.queries.DeleteOrganization(ctx, id)
}

func (s *Service) GetAllOrganizationsByUser(ctx context.Context, userID pgtype.UUID) ([]db.GetAllOrganizationsByUserRow, error) {
	return s.queries.GetAllOrganizationsByUser(ctx, userID)
}

func (s *Service) GetAllUsersInOrganization(ctx context.Context, organizationID pgtype.UUID) ([]db.GetAllUsersInOrganizationRow, error) {
	return s.queries.GetAllUsersInOrganization(ctx, organizationID)
}