package repository

import "github.com/google/uuid"

type Result struct {
	RequestID      uuid.UUID
	InputString    string
	UnpackedResult string
}

type Repository interface {
	InsertResult(requestID uuid.UUID, input string, result string) error
	SelectByID(id uuid.UUID) ([]Result, error)
}