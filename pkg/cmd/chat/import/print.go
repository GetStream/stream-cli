package imports

import (
	"errors"
	"io"
	"text/tabwriter"

	"github.com/cheynewallace/tabby"

	"github.com/GetStream/stream-cli/pkg/cmd/chat/import/validator"
)

func NewTable(w io.Writer) *tabby.Tabby {
	return tabby.NewCustom(tabwriter.NewWriter(w, 0, 0, 2, ' ', 0))
}

func PrintValidationResults(w io.Writer, results *validator.Results) {
	table := NewTable(w)

	if results.HasErrors() {
		table.AddHeader("Error", "Offset")
		for _, err := range results.Errors {
			var itemErr *validator.ItemError
			if errors.As(err, &itemErr) {
				table.AddLine(err, itemErr.Offset())
			} else {
				table.AddLine(err, "")
			}
		}
	} else {
		table.AddHeader("Type", "Count")
		for t, c := range results.Stats {
			table.AddLine(t, c)
		}
	}

	table.Print()
}
