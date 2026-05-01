package main

import (
	"os"

	"github.com/nicolaiort/shikai/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
