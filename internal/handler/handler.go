package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h *Handler) Pack(c *gin.Context) {
	var req PackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSON(c, http.StatusBadRequest, ErrorResponse{
			Error: "invalid json body",
		})
		return
	}
	res := h.unpackPrv.Pack(req.Input)
	writeJSON(c, http.StatusOK, PackResponse{
		Result: res,
	})
}

func (h *Handler) Unpack(c *gin.Context) {
	var req UnpackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSON(c, http.StatusBadRequest, ErrorResponse{
			Error: "invalid json body",
		})
		return
	}
	resp, err := h.unpackPrv.UnpackAndSave(c.Request.Context(), &unpack.UnpackAndSaveReq{
		SrcStr: req.Input,
	})
	if err != nil {
		writeJSON(c, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	writeJSON(c, http.StatusOK, UnpackResponse{
		RequestID: resp.RequestID.String(),
		Result:    resp.ResStr,
	})
}

func (h *Handler) Results(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		writeJSON(c, http.StatusBadRequest, ErrorResponse{
			Error: "id query param is required",
		})
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(c, http.StatusBadRequest, ErrorResponse{
			Error: "invalid uuid",
		})
		return
	}
	results, err := h.unpackPrv.GetByID(c.Request.Context(), id)
	if err != nil {
		writeJSON(c, http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	writeJSON(c, http.StatusOK, results)
}

func writeJSON(c *gin.Context, statusCode int, data any) {
	respBytes, err := json.Marshal(data)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to marshal response\n")
		return
	}
	c.Header("Content-Type", "application/json")
	c.Status(statusCode)
	if _, err = c.Writer.Write(append(respBytes, '\n')); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
