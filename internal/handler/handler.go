package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/timur902/strings/internal/unpack"
)

type Handler struct {
	unpackPrv *unpack.Provider
}

func NewHandler(unpackPrv *unpack.Provider) *Handler {
	return &Handler{
		unpackPrv: unpackPrv,
	}
}

func (h *Handler) Pack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, unpack.ErrorResponse{
			Error: "method not allowed",
		})
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, unpack.ErrorResponse{
			Error: "failed to read request body",
		})
		return
	}
	var req unpack.PackHTTPRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, unpack.ErrorResponse{
			Error: "invalid json body",
		})
		return
	}
	res := h.unpackPrv.Pack(req.Input)
	writeJSON(w, http.StatusOK, unpack.PackHTTPResponse{
		Result: res,
	})
}

func (h *Handler) Unpack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, unpack.ErrorResponse{
			Error: "method not allowed",
		})
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, unpack.ErrorResponse{
			Error: "failed to read request body",
		})
		return
	}
	var req unpack.UnpackHTTPRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, unpack.ErrorResponse{
			Error: "invalid json body",
		})
		return
	}
	resp, err := h.unpackPrv.UnpackAndSave(r.Context(), &unpack.UnpackAndSaveReq{
		SrcStr: req.Input,
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, unpack.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, unpack.UnpackHTTPResponse{
		RequestID: resp.RequestID.String(),
		Result:    resp.ResStr,
	})
}

func (h *Handler) Results(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, unpack.ErrorResponse{
			Error: "method not allowed",
		})
		return
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, unpack.ErrorResponse{
			Error: "id query param is required",
		})
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, unpack.ErrorResponse{
			Error: "invalid uuid",
		})
		return
	}
	results, err := h.unpackPrv.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, unpack.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, results)
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	respBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(append(respBytes, '\n'))
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
