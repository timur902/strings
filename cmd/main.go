package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/timur902/strings/internal/repository"
	"github.com/timur902/strings/internal/unpack"
)

func main() {
	input := flag.String("input", "", "input string")
	packMode := flag.Bool("pack", false, "pack string")
	unpackMode := flag.Bool("unpack", false, "unpack string")
	daemon := flag.Bool("daemon", false, "daemon mode")
	getMode := flag.String("get", "", "get results by request id")
	flag.Parse()
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
	if *daemon {
		reader := bufio.NewScanner(os.Stdin)
		for {
			select {
			case <-ctx.Done():
				log.Println("shutting down")
				return
			default:
			}
			fmt.Print("Введите строку для распаковки: ")
			if !reader.Scan() {
				return
			}
			s := reader.Text()
			resp, err := unpackPrv.UnpackAndSave(ctx, &unpack.UnpackAndSaveReq{
				SrcStr: s,
			})
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("request id:", resp.RequestID)
			fmt.Println("result:", resp.ResStr)
		}
	}
	run(ctx, unpackPrv, *input, *packMode, *unpackMode, *getMode)
}

func run(ctx context.Context, unpackPrv *unpack.Provider, s string, packMode bool, unpackMode bool, getMode string) {
	if packMode {
		fmt.Println(unpackPrv.Pack(s))
		return
	}
	if unpackMode {
		resp, err := unpackPrv.UnpackAndSave(ctx, &unpack.UnpackAndSaveReq{
			SrcStr: s,
		})
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("request id:", resp.RequestID)
		fmt.Println("result:", resp.ResStr)
		return
	}
	if getMode != "" {
		id, err := uuid.Parse(getMode)
		if err != nil {
			log.Println("invalid uuid:", err)
			return
		}
		results, err := unpackPrv.GetByID(ctx, id)
		if err != nil {
			log.Println(err)
			return
		}
		if len(results) == 0 {
			log.Println("no results found")
			return
		}
		for _, res := range results {
			fmt.Println("request id:", res.RequestID)
			fmt.Println("input string:", res.InputString)
			fmt.Println("unpacked result:", res.UnpackedResult)
			fmt.Println("---")
		}
		return
	}
	log.Println("specify --pack or --unpack or --get")
}