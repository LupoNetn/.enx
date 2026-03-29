package env

import (
	"context"

	"github.com/google/uuid"
	"github.com/luponetn/enx/internal/db"
)

type Svc interface {
	CreateEnv(ctx context.Context, args db.CreateEnvParams) (db.Envs, error)
	GetEnvByID(ctx context.Context, id uuid.UUID) (db.Envs, error)
	GetEnvsByProject(ctx context.Context, projectID uuid.UUID) ([]db.Envs, error)
	GetEnvByNameInProject(ctx context.Context, args db.GetEnvByNameInProjectParams) (db.Envs, error)
	UpdateEnv(ctx context.Context, args db.UpdateEnvParams) (db.Envs, error)
	DeleteEnv(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) Svc {
	return &Service{queries: queries}
}

func (s *Service) CreateEnv(ctx context.Context, args db.CreateEnvParams) (db.Envs, error) {
	return s.queries.CreateEnv(ctx, args)
}

func (s *Service) GetEnvByID(ctx context.Context, id uuid.UUID) (db.Envs, error) {
	return s.queries.GetEnvByID(ctx, id)
}

func (s *Service) GetEnvsByProject(ctx context.Context, projectID uuid.UUID) ([]db.Envs, error) {
	return s.queries.GetEnvsByProject(ctx, projectID)
}

func (s *Service) GetEnvByNameInProject(ctx context.Context, args db.GetEnvByNameInProjectParams) (db.Envs, error) {
	return s.queries.GetEnvByNameInProject(ctx, args)
}

func (s *Service) UpdateEnv(ctx context.Context, args db.UpdateEnvParams) (db.Envs, error) {
	return s.queries.UpdateEnv(ctx, args)
}

func (s *Service) DeleteEnv(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeleteEnv(ctx, id)
}
