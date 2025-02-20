package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var usage = `
A cli programme to make or compare bcrypt passwords

bcrypt has a string:
<programme> "password to hash"

check if a string matches a bcrypt hash
<programme> -hash "hashvalue" "password to check"

`

func flagGet() (string, string) {

	var (
		hash string
	)

	flag.StringVar(&hash, "hash", "", "hash value to match")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(flag.CommandLine.Output(), usage)
	}

	flag.Parse()
	passwords := flag.Args()

	if hash != "" && len(hash) < 59 { // bcrypt.ErrHashTooShort
		fmt.Println("hash is too short. Maybe it needs to be quoted '....'?")
		os.Exit(1)
	}

	if len(passwords) != 1 {
		fmt.Println("please provide a single password")
		flag.Usage()
		os.Exit(1)
	}

	return hash, passwords[0]
}

func main() {
	hash, password := flagGet()

	// match mode
	if hash != "" {
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err != nil {
			fmt.Printf("match error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("match!")
		return
	}

	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost+2)
	if err != nil {
		fmt.Printf("hashing error: %v\n", err)
		os.Exit(1)
	}
	// hash is b64 encoded
	fmt.Printf(string(h))
}
