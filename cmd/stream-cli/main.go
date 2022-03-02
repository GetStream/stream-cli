package main

import (
	"fmt"
	"os"

	"github.com/GetStream/stream-cli/pkg/cmd/root"
)

func main() {
	if err := mainRun(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func mainRun() error {
	return root.NewRootCmd().Execute()
}
