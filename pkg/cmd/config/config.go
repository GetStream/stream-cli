package config

import (
	"fmt"
	"text/tabwriter"

	"github.com/AlecAivazis/survey/v2"
	cfg "github.com/GetStream/stream-cli/pkg/config"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
)

func NewRootConfigCmd() *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuestionnaire(cmd)
		},
	}
}

func removeAppCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [app-name]",
		Short: "Remove an application. This operation is irrevocable",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			config := cfg.GetConfig(cmd)
			return config.Remove(args[0])
		},
	}
}

func listAppsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all applications",

		RunE: func(cmd *cobra.Command, args []string) error {
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			t := tabby.NewCustom(w)
			t.AddHeader("", "Name", "Access Key", "Secret Key", "Url")

			config := cfg.GetConfig(cmd)

			for _, app := range config.Apps {
				def := ""
				if app.Name == config.Default {
					def = "(default)"
				}

				secret := fmt.Sprintf("**************%v", app.AccessSecretKey[len(app.AccessSecretKey)-4:])
				t.AddLine(def, app.Name, app.AccessKey, secret, app.URL)
			}
			t.Print()
			return nil
		},
	}
}

func setAppDefaultCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default [app-name]",
		Short: "Set an application as the default",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			config := cfg.GetConfig(cmd)
			return config.SetDefault(args[0])
		},
	}

	return cmd
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
			Name:   "URL",
			Prompt: &survey.Input{Message: "(optional) Which base URL do you want to use? Default value is our edge URL."},
			Transform: func(ans interface{}) interface{} {
				s, ok := ans.(string)

				if !ok || s == "" {
					return cfg.DefaultEdgeURL
				}

				return s
			},
		},
	}
}
