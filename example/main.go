package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/MarkMandriota/TinyVM"
)

var vm *tinyvm.Machine

func init() {
	if len(os.Args) < 2 {
		log.Fatalf("no input file")
	}

	vm = new(tinyvm.Machine)
	vm.Init(nil)

	fi, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("error while opening file: %v", err)
	}
	defer fi.Close()

	vm.Text, _ = io.ReadAll(fi)
}

func main() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)

	defer func() {
		w.Flush()

		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	vm.Execute(r, w)
}
