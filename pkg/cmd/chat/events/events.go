package events

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
	"github.com/GetStream/stream-cli/pkg/version"
	"github.com/MakeNowJust/heredoc"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{listenCmd()}
}

func listenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listen-events --user-id [user-id]",
		Short: "Listen to events",
		Long: heredoc.Doc(`
			The command opens a WebSocket connection to the backend in the name of the user
			and prints the received events to the standard output.
			Press Ctrl+C to exit.
		`),
		Example: heredoc.Doc(`
			# Listen to events for user with id 'my-user-1'
			$ stream-cli chat listen-events --user-id my-user-1
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			userID, _ := cmd.Flags().GetString("user-id")
			config := config.GetConfig(cmd)
			client, _ := config.GetClient(cmd)
			token, err := client.CreateToken(userID, time.Time{})
			if err != nil {
				return err
			}

			apiKey, _, err := config.GetCredentials(cmd)
			if err != nil {
				return err
			}

			cmd.Println("> 🚨 Warning! The WebSocket connection can be expensive so we close it down after 60 seconds.")
			time.Sleep(2 * time.Second)
			// Giving the user 2 seconds to read the warning message
			// because the first heartbeat is sent super quickly and
			// takes up the whole screen.

			url := getUrl(userID, apiKey, token)
			websocket.DefaultDialer.HandshakeTimeout = time.Second * 5
			conn, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				return err
			}
			cmd.Println("> Successfully connected. Waiting for events...⌛️")

			exit := make(chan os.Signal, 1)
			signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

			// Since keeping connections can be expensive
			// let's just exit after 60 seconds.
			go func() {
				time.Sleep(50 * time.Second)
				cmd.Println("> Exiting in 10 seconds...")
				time.Sleep(10 * time.Second)
				cmd.Println("> 60 seconds passed. Exiting now.")
				exit <- syscall.SIGINT
			}()

			go func() {
				for {
					var event map[string]any
					err := conn.ReadJSON(&event)
					if err != nil {
						cmd.PrintErr(err)
						continue
					}

					err = utils.PrintObject(cmd, event)
					if err != nil {
						cmd.PrintErr(err)
						continue
					}

					cmd.Println("> Press Ctrl+C to exit")
				}
			}()

			<-exit
			cmd.Println("> Exiting")
			return conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		},
	}

	fl := cmd.Flags()
	fl.StringP("user-id", "u", "", "User ID")
	fl.StringP("output-format", "o", "json", "Output format")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}

func getUrl(userID, apiKey, token string) string {
	u, _ := url.Parse("wss://chat.stream-io-api.com/connect")

	payload := map[string]any{
		"user_id":      userID,
		"user_details": map[string]string{"id": userID},
	}
	b, _ := json.Marshal(payload)

	params := url.Values{}
	params.Add("json", string(b))
	params.Add("api_key", apiKey)
	params.Add("authorization", token)
	params.Add("stream-auth-type", "jwt")
	params.Add("X-Stream-Client", fmt.Sprintf("stream-cli-%s", version.FmtVersion()))
	u.RawQuery = params.Encode()

	return u.String()
}
