package main

import (
	"flag"
	"fmt"
	"os"

	uuid "github.com/satori/go.uuid"
)

func main() {

	n := flag.Int("n", 1, "Number of UUIDs to print.")

	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	for i := 0; i < *n; i++ {
		u1 := uuid.NewV4()
		fmt.Printf("%s\n", u1)
	}

}
