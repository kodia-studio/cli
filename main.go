package main

import (
	"fmt"
	"os"

	"github.com/kodia-studio/kodia/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
