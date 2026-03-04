package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/timur902/stack_queue/internal/unpack"
)

func main() {
	input := flag.String("input", "", "input string")
	packMode := flag.Bool("pack", false, "pack string")
	unpackMode := flag.Bool("unpack", false, "unpack string")
	daemon := flag.Bool("daemon", false, "daemon mode")

	// Создаёшь NewRepository в том же пакете что и интерфейс, возвращаешь структуру которая его реализует
	// Пробрасываешь в unpackPrv и вызываешь метода для бд уже оттуда.
	// И соответсвенно всю логику работы с базой данных ты пишешь в методах структуры репозитория которую создал

	unpackPrv := unpack.NewProvider()

	flag.Parse()
	if *daemon {
		reader := bufio.NewScanner(os.Stdin)
		for {

			fmt.Print("Введите строку: ")

			if !reader.Scan() {
				return
			}

			s := reader.Text()

			run(unpackPrv, s, *packMode, *unpackMode)
		}
	}
	run(unpackPrv, *input, *packMode, *unpackMode)
}

func run(unpackPrv *unpack.Provider, s string, packMode bool, unpackMode bool) {
	if packMode {

		fmt.Println(unpackPrv.Pack(s))
		return
	}
	if unpackMode {

		res, err := unpackPrv.Unpack(s)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(res)
		return
	}

	log.Println("specify --pack or --unpack")
}
