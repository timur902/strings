package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/timur902/strings/internal/handler"
	"github.com/timur902/strings/internal/repository"
	"github.com/timur902/strings/internal/unpack"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("warning: failed to load .env: %v", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		log.Printf("warning: failed to check .env: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}
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
	h := handler.NewHandler(unpackPrv)
	router := gin.Default()

	router.POST("/pack", h.Pack)
	router.POST("/unpack", h.Unpack)
	router.GET("/results", h.Results)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
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
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
