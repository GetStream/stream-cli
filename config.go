package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	configDir  = ".stream-cli"
	configFile = "config"

	defaultEdgeURL = "https://chat.stream-io-api.com"
)

func NewRootConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "config [options] <command> <subcommand>",
		Long: `Manage app configurations`,
		Run:  func(_ *cobra.Command, _ []string) {},
	}

	cmd.AddCommand(newConfigCmd())

	return cmd
}

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "new",
		Long:          "Create a new app configuration",
		SilenceUsage:  true,
		SilenceErrors: true,

		RunE: func(_ *cobra.Command, _ []string) error {
			d, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("cannot get user's home directory: %v", err)
			}

			err = os.Mkdir(path.Join(d, configDir), 0755)
			if err != nil && !os.IsExist(err) {
				return fmt.Errorf("cannot create config directory: %v", err)
			}

			filepath := path.Join(d, configDir, configFile)

			f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
			defer f.Close()

			newConfig := newDefaultConfig()
			err = survey.Ask(questions(), &newConfig)
			if err != nil {
				return err
			}

			return addNewConfig(f, &newConfig)
		},
	}

	return cmd
}

type appConfig struct {
	Name            string `yaml:"-"`
	AccessKey       string `yaml:"access-key"`
	AccessSecretKey string `yaml:"access-secret-key"`
	URL             string `yaml:"url"`
}

func newDefaultConfig() appConfig {
	return appConfig{
		URL: defaultEdgeURL,
	}
}

// addNewConfig adds a new app configuration.
func addNewConfig(file *os.File, newConfig *appConfig) error {
	file.Seek(0, io.SeekStart)
	appsConfig := make(map[string]*appConfig)

	err := yaml.NewDecoder(file).Decode(appsConfig)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	if _, ok := appsConfig[newConfig.Name]; ok {
		return fmt.Errorf("configuration for %q already exists", newConfig.Name)
	}

	newCfg := map[string]*appConfig{
		newConfig.Name: newConfig,
	}
	err = yaml.NewEncoder(file).Encode(newCfg)
	return err
}

// questions returns all questions to ask to configure an app.
func questions() []*survey.Question {
	return []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "What is the name of your app? (eg. prod, staging, testing)"},
			Validate: survey.Required,
		},
		{
			Name:     "accessKey",
			Prompt:   &survey.Input{Message: "What is your access key?"},
			Validate: survey.Required,
		},
		{
			Name:     "accessSecretKey",
			Prompt:   &survey.Password{Message: "What is your access secret key?"},
			Validate: survey.Required,
		},
		{
			Name:   "URL",
			Prompt: &survey.Input{Message: "(optional) Which base URL do you want to use? Default value is our edge URL."},
		},
	}
}
