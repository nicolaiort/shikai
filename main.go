package main

import (
	"os"

	"github.com/shikai/release/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
