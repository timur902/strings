package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/timur902/strings/internal/handler"
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
	handler := handler.NewHandler(unpackPrv)
	mux := http.NewServeMux()

	mux.HandleFunc("/pack", handler.Pack)
	mux.HandleFunc("/unpack", handler.Unpack)
	mux.HandleFunc("/results", handler.Results)

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
