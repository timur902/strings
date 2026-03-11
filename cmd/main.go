package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
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
	dbURL := "postgres://postgres:postgres@localhost:5432/unpacker"
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("failed to create db pool: %v", err)
	}
	defer pool.Close()
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}
	repo := repository.NewRepository(pool)
	unpackPrv := unpack.NewProvider(repo)
	if *daemon {
		reader := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Введите строку для распаковки: ")
			if !reader.Scan() {
				return
			}
			s := reader.Text()
			requestID, res, err := unpackPrv.UnpackAndSave(s)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("request id:", requestID)
			fmt.Println("result:", res)
		}
	}
	run(unpackPrv, *input, *packMode, *unpackMode, *getMode)
}

func run(unpackPrv *unpack.Provider, s string, packMode bool, unpackMode bool, getMode string) {
	if packMode {
		fmt.Println(unpackPrv.Pack(s))
		return
	}
	if unpackMode {
		requestID, res, err := unpackPrv.UnpackAndSave(s)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("request id:", requestID)
		fmt.Println("result:", res)
		return
	}
	if getMode != "" {
		id, err := uuid.Parse(getMode)
		if err != nil {
			log.Println("invalid uuid:", err)
			return
		}
		results, err := unpackPrv.GetByID(id)
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