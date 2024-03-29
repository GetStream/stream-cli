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
		Use:   "listen-events --user-id [user-id] --timeout [number]",
		Short: "Listen to events",
		Long: heredoc.Doc(`
			The command opens a WebSocket connection to the backend in the name of the user
			and prints the received events to the standard output.
			Press Ctrl+C to exit.
		`),
		Example: heredoc.Doc(`
			# Listen to events for user with id 'my-user-1'
			$ stream-cli chat listen-events --user-id my-user-1

			# Listen to events for user with id 'my-user-2' and keeping the connection open for 120 seconds
			$ stream-cli chat listen-events --user-id my-user-1 --timeout 120
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			timeout, _ := cmd.Flags().GetInt32("timeout")
			if timeout > 300 {
				return fmt.Errorf("timeout cannot be greater than 300")
			}
			userID, _ := cmd.Flags().GetString("user-id")
			config := config.GetConfig(cmd)

			app, err := config.GetDefaultAppOrExplicit(cmd)
			if err != nil {
				return err
			}

			client, err := config.GetClient(cmd)
			if err != nil {
				return err
			}

			token, err := client.CreateToken(userID, time.Time{})
			if err != nil {
				return err
			}

			cmd.Printf("> 🚨 Warning! The WebSocket connection can be expensive so we close it after %d seconds.\n", timeout)
			time.Sleep(2 * time.Second)
			// Giving the user 2 seconds to read the warning message
			// because the first heartbeat is sent super quickly and
			// takes up the whole screen.

			url, err := getUrl(app, userID, token)
			if err != nil {
				return err
			}

			websocket.DefaultDialer.HandshakeTimeout = time.Second * 5
			conn, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				return err
			}
			cmd.Println("> Successfully connected. Waiting for events...⌛️")

			exit := make(chan os.Signal, 1)
			signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

			// Since keeping connections can be expensive
			// let's just exit after 'timeout' seconds.
			go func() {
				time.Sleep(time.Duration(timeout-10) * time.Second)
				cmd.Println("> Exiting in 10 seconds...")
				time.Sleep(10 * time.Second)
				cmd.Printf("> %d seconds passed. Exiting now.\n", timeout)
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
	fl.Int32P("timeout", "t", 60, "[optional] For how many seconds do we keep the connection alive. Default is 60 seconds, max is 300.")
	fl.StringP("user-id", "u", "", "[required] User ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}

func getUrl(app *config.App, userID, token string) (string, error) {
	u, err := url.Parse(app.ChatURL)
	if err != nil {
		return "", err
	}

	u.Scheme = "wss"
	u.Path = "connect"

	payload := map[string]any{
		"user_id":      userID,
		"user_details": map[string]string{"id": userID},
	}
	b, _ := json.Marshal(payload)

	params := url.Values{}
	params.Add("json", string(b))
	params.Add("api_key", app.AccessKey)
	params.Add("authorization", token)
	params.Add("stream-auth-type", "jwt")
	params.Add("X-Stream-Client", fmt.Sprintf("stream-cli-%s", version.FmtVersion()))
	u.RawQuery = params.Encode()

	return u.String(), nil
}
