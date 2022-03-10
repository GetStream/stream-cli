package config

import (
	"errors"
	"fmt"
	"net/url"
	"text/tabwriter"

	"github.com/AlecAivazis/survey/v2"
	cfg "github.com/GetStream/stream-cli/pkg/config"
	"github.com/MakeNowJust/heredoc"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage app configurations",
	}

	cmd.AddCommand(newAppCmd(), removeAppCmd(), listAppsCmd(), setAppDefaultCmd())

	return cmd
}

func newAppCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "Add a new application",
		Long:  "Add a new application which can be used for further operations",
		Example: heredoc.Doc(`
			# Add a new application to the CLI
			$ stream-cli config new
			? What is the name of your app? (eg. prod, staging, testing) testing
			? What is your access key? abcd1234efgh456
			? What is your access secret key? ***********************************
			? (optional) Which base URL do you want to use for Chat? https://chat.stream-io-api.com

			Application successfully added. 🚀
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuestionnaire(cmd)
		},
	}
}

func removeAppCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [app-name-1] [app-name-2] [app-name-n]",
		Short: "Remove one or more application.",
		Long:  "Remove one or more application from the configuraiton file. This operation is irrevocable.",
		Example: heredoc.Doc(`
			# Remove a single application from the CLI
			$ stream-cli config remove staging

			# Remove multiple applications from the CLI
			$ stream-cli config remove staging testing
		`),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := cfg.GetConfig(cmd)
			for _, appName := range args {
				if err := config.Remove(appName); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func listAppsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all applications",
		Long:  "List all applications which are configured in the configuration file",
		Example: heredoc.Doc(`
			# List all applications
			$ stream-cli config list
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			t := tabby.NewCustom(w)
			t.AddHeader("", "Name", "Access Key", "Secret Key", "URL")

			config := cfg.GetConfig(cmd)

			for _, app := range config.Apps {
				def := ""
				if app.Name == config.Default {
					def = "(default)"
				}

				secret := fmt.Sprintf("**************%v", app.AccessSecretKey[len(app.AccessSecretKey)-4:])
				t.AddLine(def, app.Name, app.AccessKey, secret, app.ChatURL)
			}
			t.Print()
			return nil
		},
	}
}

func setAppDefaultCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "default [app-name]",
		Short: "Set an application as the default",
		Long: heredoc.Doc(`
			Set an application as the default which will be used
			for all further operations unless specified otherwise.
		`),
		Example: heredoc.Doc(`
			# Set an application as the default
			$ stream-cli config default staging

			# All underlying operations will use it if not specified otherwise
			$ stream-cli chat get-app
			# Prints the settings of staging app

			# Specifying other apps during an operation
			$ stream-cli chat get-app --app prod
			# Prints the settings of prod app
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := cfg.GetConfig(cmd)
			return config.SetDefault(args[0])
		},
	}
}

func runQuestionnaire(cmd *cobra.Command) error {
	var newAppConfig cfg.App
	err := survey.Ask(questions(), &newAppConfig)
	if err != nil {
		return err
	}

	config := cfg.GetConfig(cmd)
	err = config.Add(newAppConfig)
	if err != nil {
		return err
	}

	cmd.Println("Application successfully added. 🚀")
	return nil
}

func questions() []*survey.Question {
	return []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "What is the name of your app? (eg. prod, staging, testing)"},
			Validate: survey.Required,
		},
		{
			Name:   "accessKey",
			Prompt: &survey.Input{Message: "What is your access key?"},
			Validate: survey.ComposeValidators(
				survey.Required,
				survey.MinLength(10),
				survey.MaxLength(30)),
		},
		{
			Name:   "accessSecretKey",
			Prompt: &survey.Password{Message: "What is your access secret key?"},
			Validate: survey.ComposeValidators(
				survey.Required,
				survey.MinLength(50),
				survey.MaxLength(75)),
		},
		{
			Name: "ChatURL",
			Prompt: &survey.Input{
				Message: "(optional) Which base URL do you want to use for Chat?",
				Default: cfg.DefaultChatEdgeURL,
			},
			Validate: func(ans interface{}) error {
				u, ok := ans.(string)
				if !ok {
					return errors.New("invalid url")
				}

				_, err := url.ParseRequestURI(u)
				if err != nil {
					return errors.New("invalid url format. make sure it matches <scheme>://<host>")
				}
				return nil
			},
			Transform: func(ans interface{}) interface{} {
				s, ok := ans.(string)

				if !ok || s == "" {
					return cfg.DefaultChatEdgeURL
				}

				return s
			},
		},
	}
}
