package repository

import "github.com/google/uuid"

type InsertResultReq struct {
	RequestID      uuid.UUID
	InputString    string
	UnpackedResult string
}

type Result struct {
	RequestID      uuid.UUID `json:"request_id"`
	InputString    string    `json:"input_string"`
	UnpackedResult string    `json:"unpacked_result"`
}