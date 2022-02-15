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
	config, err := cli.NewConfig()
	if err != nil {
		return err
	}
	defer config.Close()
	rootCmd := cli.NewRootCmd(config)
	err = rootCmd.Run(os.Args)
	return err
}
