package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/MarkMandriota/tinyvm"
)

func main() {
	flag.Parse()

	fs, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatalf("error while opening file \"%s\": %v", flag.Arg(0), err)
	}
	defer fs.Close()

	vm := tinyvm.NewMachine(nil)
	vm.Text, _ = io.ReadAll(fs)

	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	beg := time.Now()

	if err := vm.Execute(r, w); err != nil && err != io.EOF {
		log.Println(err)
	}

	log.Printf("Total execution time: %v", time.Since(beg))
}
