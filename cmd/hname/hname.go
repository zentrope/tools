package main

import (
	"fmt"
	"os"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Unable to figure out the hostname: [%v]\n", err)
	} else {
		fmt.Println(hostname)
	}
}
