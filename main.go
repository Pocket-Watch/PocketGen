package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello, World")

	exec, err := os.Executable()
	if err != nil {
		exec = ""
	}
	fmt.Println("Usage:")
	fmt.Printf("  %v\n", exec)
}
