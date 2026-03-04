package main

import (
	"bufio"
	"go_test1/task3/unpack"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	input := flag.String("input", "", "input string")
	packMode := flag.Bool("pack", false, "pack string")
	unpackMode := flag.Bool("unpack", false, "unpack string")
	daemon := flag.Bool("daemon", false, "daemon mode")
	flag.Parse()
	if *daemon {
		reader := bufio.NewScanner(os.Stdin)
		for {

			fmt.Print("Введите строку: ")

			if !reader.Scan() {
				return
			}

			s := reader.Text()

			run(s, *packMode, *unpackMode)
		}
	}
	run(*input, *packMode, *unpackMode)
}
func run(s string, packMode bool, unpackMode bool) {
	if packMode {

		fmt.Println(unpack.Pack(s))
		return
	}
	if unpackMode {

		res, err := unpack.Unpack(s)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(res)
		return
	}

	log.Println("specify --pack or --unpack")
}