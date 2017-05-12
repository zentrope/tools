package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

const DEFAULT = "password"

func main() {
	e := flag.String("e", DEFAULT, "Password to hash or compare")
	d := flag.String("d", DEFAULT, "Hash to compare with password")

	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	if (*e != DEFAULT) && (*d == DEFAULT) {
		fmt.Printf("Hashing [%s].\n", *e)
		result, err := bcrypt.GenerateFromPassword([]byte(*e), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%x\v", result)
	}

	if *d != DEFAULT {
		decoded, err := hex.DecodeString(*d)
		if err != nil {
			panic(err)
		}

		err = bcrypt.CompareHashAndPassword(decoded, []byte(*e))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Password is a match.")
		}
	}

}
