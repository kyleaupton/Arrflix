package main

import (
	"fmt"
	"os"

	"github.com/kyleaupton/arrflix/internal/password"
)

func main() {
	// read password from first argument
	raw := os.Args[1]

	hashed, err := password.Hash(raw)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}
	fmt.Println(hashed)
}
