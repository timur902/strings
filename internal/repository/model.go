package repository

import "github.com/google/uuid"

type InsertResultReq struct {
	RequestID      uuid.UUID
	InputString    string
	UnpackedResult string
}

type Result struct {
	RequestID      uuid.UUID
	InputString    string
	UnpackedResult string
}