package repository

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	InsertResult(ctx context.Context, req *InsertResultReq) error
	SelectByID(ctx context.Context, id uuid.UUID) ([]Result, error)
}