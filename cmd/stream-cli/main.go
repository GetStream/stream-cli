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
	d, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("cannot get user's home directory: %v", err)
	}

	config, err := cli.NewConfig(d)
	if err != nil {
		return err
	}
	rootCmd := cli.NewRootCmd(config)
	return rootCmd.Run(os.Args)
}
