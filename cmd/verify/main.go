package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run verify.go <password> <hash>")
		os.Exit(1)
	}

	password := os.Args[1]
	hash := os.Args[2]

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Password does NOT match:", err)
		os.Exit(1)
	}

	fmt.Println("Password matches!")
}
