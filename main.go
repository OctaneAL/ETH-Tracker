package main

import (
	"os"

	"github.com/OctaneAL/ETH-Tracker/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
