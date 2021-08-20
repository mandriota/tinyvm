package main

import (
	"io"
	"log"
	"os"
	"time"

	tinyvm "github.com/MarkMandriota/TinyVM"
)

var vm *tinyvm.Machine

func init() {
	if len(os.Args) < 2 {
		log.Fatalf("no input file")
	}

	fi, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("error while opening file: %v", err)
	}
	defer fi.Close()

	vm = tinyvm.NewMachine(nil, os.Stdin, os.Stdout)
	vm.Text, _ = io.ReadAll(fi)
}

func main() {
	beg := time.Now()

	defer func() {
		log.Printf("Total execution time: %v", time.Since(beg))

		if v := recover(); v != nil {
			log.Printf("MSG!: %v", v)
		}
	}()

	vm.Execute()
}
