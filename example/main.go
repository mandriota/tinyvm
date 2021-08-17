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
		log.Fatalf("error while opening file: ", err)
	}
	defer fi.Close()

	vm.Text, _ = io.ReadAll(fi)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	vm.Execute(r, w)
}
