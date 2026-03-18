package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/timur902/strings/internal/repository"
	"github.com/timur902/strings/internal/unpack"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	dbURL := "postgres://postgres:postgres@localhost:5432/unpacker"
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to create db pool: %v", err)
	}
	defer pool.Close()
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}
	repo := repository.NewRepository(pool)
	unpackPrv := unpack.NewProvider(repo)
	mux := http.NewServeMux()

	mux.HandleFunc("/pack", func(w http.ResponseWriter, r *http.Request) {
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
		res := unpackPrv.Pack(req.Input)
		writeJSON(w, http.StatusOK, unpack.PackHTTPResponse{
			Result: res,
		})
	})

	mux.HandleFunc("/unpack", func(w http.ResponseWriter, r *http.Request) {
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
		resp, err := unpackPrv.UnpackAndSave(r.Context(), &unpack.UnpackAndSaveReq{
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
	})

	mux.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
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
		results, err := unpackPrv.GetByID(r.Context(), id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, unpack.ErrorResponse{
				Error: err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusOK, results)
	})
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		<-ctx.Done()
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()
	log.Println("http server started on :8080")
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
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