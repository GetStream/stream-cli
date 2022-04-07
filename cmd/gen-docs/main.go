package main

import (
	"fmt"
	"os"

	"github.com/GetStream/stream-cli/pkg/cmd/root"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fl := pflag.FlagSet{}
	dir := fl.StringP("output", "o", "./docs", "Path directory where you want generate doc files")

	if err := fl.Parse(args); err != nil {
		return err
	}

	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		if err := os.MkdirAll(*dir, 0o755); err != nil {
			return err
		}
	}

	rootCmd := root.NewCmd()
	rootCmd.DisableAutoGenTag = true
	return doc.GenMarkdownTree(rootCmd, *dir)
}
