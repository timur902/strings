package unpack

import "github.com/google/uuid"

type UnpackAndSaveReq struct {
	SrcStr string
}

type UnpackAndSaveResp struct {
	RequestID uuid.UUID
	ResStr    string
}