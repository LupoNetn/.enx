package env

import "github.com/luponetn/enx/internal/db"

type Svc interface {
}

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) Svc {
	return &Service{queries: queries}
}


//implement handlers for creating of env variables
