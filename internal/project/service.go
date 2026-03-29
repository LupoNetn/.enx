package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/luponetn/enx/internal/db"
)

type Service struct {
	queries *db.Queries
}

type Svc interface {
	CreateProject(ctx context.Context, args db.CreateProjectParams) (db.CreateProjectRow, error)
	UpdateProject(ctx context.Context, args db.UpdateProjectParams) (db.UpdateProjectRow, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error
	GetProjectsByUser(ctx context.Context, userID uuid.UUID) ([]db.GetProjectsByUserRow, error)
	GetAllUsersInProject(ctx context.Context, projectID uuid.UUID) ([]db.GetAllUsersInProjectRow, error)
	GetProjectByName(ctx context.Context, args db.GetProjectByNameParams) (db.GetProjectByNameRow, error)
	GetProjectByNameForUser(ctx context.Context, args db.GetProjectByNameForUserParams) (db.GetProjectByNameForUserRow, error)
}

func NewProjectService(queries *db.Queries) Svc {
	return &Service{
		queries: queries,
	}
}

func (s *Service) CreateProject(ctx context.Context, args db.CreateProjectParams) (db.CreateProjectRow, error) {
	return s.queries.CreateProject(ctx, args)
}

func (s *Service) UpdateProject(ctx context.Context, args db.UpdateProjectParams) (db.UpdateProjectRow, error) {
	return s.queries.UpdateProject(ctx, args)
}

func (s *Service) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeleteProject(ctx, id)
}

func (s *Service) GetProjectsByUser(ctx context.Context, userID uuid.UUID) ([]db.GetProjectsByUserRow, error) {
	return s.queries.GetProjectsByUser(ctx, userID)
}

func (s *Service) GetAllUsersInProject(ctx context.Context, projectID uuid.UUID) ([]db.GetAllUsersInProjectRow, error) {
	return s.queries.GetAllUsersInProject(ctx, projectID)
}

func (s *Service) GetProjectByName(ctx context.Context, args db.GetProjectByNameParams) (db.GetProjectByNameRow, error) {
	return s.queries.GetProjectByName(ctx, args)
}

func (s *Service) GetProjectByNameForUser(ctx context.Context, args db.GetProjectByNameForUserParams) (db.GetProjectByNameForUserRow, error) {
	return s.queries.GetProjectByNameForUser(ctx, args)
}
