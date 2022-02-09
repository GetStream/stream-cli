package main

import (
	"fmt"
	"os"

	"github.com/GetStream/stream-cli"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	rootCmd := cli.NewRootCmd()
	err := rootCmd.Run(os.Args)
	return err
}
