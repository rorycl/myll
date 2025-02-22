package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var usage = `
A cli programme to make or compare bcrypt passwords

bcrypt has a string:
<programme> "password to hash"

check if a string matches a bcrypt hash
<programme> -hash "hashvalue" "password to check"

`

func flagGet() (string, string, int, bool) {

	var (
		hash    string
		cost    int
		verbose bool
	)

	flag.StringVar(&hash, "hash", "", "hash value to match")
	flag.IntVar(&cost, "cost", bcrypt.DefaultCost, "bcrypt cost")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode")

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

	return hash, passwords[0], cost, verbose
}

func main() {
	hash, password, cost, verbose := flagGet()

	t := time.Now()

	// match mode
	if hash != "" {
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err != nil {
			fmt.Printf("match error: %v\n", err)
			os.Exit(1)
		}
		if verbose {
			fmt.Println("comparison duration:", time.Now().Sub(t))
		}
		fmt.Println("match ok")
		return
	}

	if verbose {
		fmt.Printf("using cost %d (default %d)\n", cost, bcrypt.DefaultCost)
	}
	h, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		fmt.Printf("hashing error: %v\n", err)
		os.Exit(1)
	}
	if verbose {
		fmt.Println("generation duration:", time.Now().Sub(t))
	}
	// hash is b64 encoded
	fmt.Println(string(h))
}
