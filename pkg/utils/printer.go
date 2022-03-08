package utils

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func PrintObject(cmd *cobra.Command, object interface{}) error {
	format, err := cmd.Flags().GetString("output-format")
	if err != nil {
		return err
	}

	switch format {
	case "json":
		return printJsonObject(cmd, object)
	default:
		return fmt.Errorf("unknown output format: %s", format)
	}
}

func printJsonObject(cmd *cobra.Command, object interface{}) error {
	var b bytes.Buffer
	jsonEncoder := json.NewEncoder(&b)
	jsonEncoder.SetIndent("", "  ")

	err := jsonEncoder.Encode(object)
	if err != nil {
		return err
	}

	cmd.Println(b.String())

	return nil
}
