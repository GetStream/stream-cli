package search

import (
	"encoding/json"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		searchCmd(),
	}
}

func searchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search --query [text] --filter [raw-json] --output-format [json|tree]",
		Short: "Search messages across channels",
		Long: heredoc.Doc(`
			Search messages across channels. You must provide a filter for the channels
			to search in, and a query string or message filter.
		`),
		Example: heredoc.Doc(`
			# Search for messages containing 'hello' in messaging channels
			$ stream-cli chat search --query "hello" --filter '{"type":"messaging"}'

			# Search with a message filter
			$ stream-cli chat search --filter '{"type":"messaging"}' --message-filter '{"text":{"$autocomplete":"hello"}}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			query, _ := cmd.Flags().GetString("query")
			filterStr, _ := cmd.Flags().GetString("filter")
			msgFilterStr, _ := cmd.Flags().GetString("message-filter")
			limit, _ := cmd.Flags().GetInt("limit")

			var filter map[string]interface{}
			if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
				return err
			}

			req := &stream.SearchRequest{
				Query:  query,
				Limit:  limit,
				Offset: 0,
			}

			if msgFilterStr != "" {
				var msgFilter map[string]interface{}
				if err := json.Unmarshal([]byte(msgFilterStr), &msgFilter); err != nil {
					return err
				}
				req.MessageFilters = msgFilter
			}

			resp, err := c.Search(cmd.Context(), stream.SearchRequest{
				Query:          req.Query,
				Filters:        filter,
				MessageFilters: req.MessageFilters,
				Limit:          req.Limit,
			})
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, resp)
		},
	}

	fl := cmd.Flags()
	fl.StringP("query", "q", "", "[optional] Search query text")
	fl.StringP("filter", "f", "", "[required] Channel filter conditions as JSON")
	fl.String("message-filter", "", "[optional] Message filter conditions as JSON")
	fl.IntP("limit", "l", 20, "[optional] Number of results to return")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")
	_ = cmd.MarkFlagRequired("filter")

	return cmd
}
