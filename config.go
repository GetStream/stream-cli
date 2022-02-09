package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cheynewallace/tabby"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

const (
	configDir  = ".stream-cli"
	configFile = "config"

	defaultEdgeURL = "https://chat.stream-io-api.com"
)

type appConfig struct {
	Name            string `yaml:"-"`
	AccessKey       string `yaml:"access-key"`
	AccessSecretKey string `yaml:"access-secret-key"`
	URL             string `yaml:"url"`
	Default         bool   `yaml:"default,omitempty"`
}

func newDefaultConfig() appConfig {
	return appConfig{
		URL: defaultEdgeURL,
	}
}

func NewRootConfigCmd() *cli.Command {
	cmd := &cli.Command{
		Name:        "config",
		Usage:       "Manage app configurations",
		Description: `Manage app configurations`,
		Subcommands: []*cli.Command{
			newConfigCmd(),
			removeConfigCmd(),
			listConfigsCmd(),
		},
	}

	return cmd
}

func newConfigCmd() *cli.Command {
	cmd := &cli.Command{
		Name:        "new",
		Usage:       "Add a new App configuration",
		UsageText:   "stream-cli config new",
		Description: "Add a new App configuration which can be used on further operations",

		//SilenceUsage:  true,
		//SilenceErrors: true,

		Action: func(ctx *cli.Context) error {
			f, err := getConfigurationFile()
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

// addNewConfig adds a new app configuration.
func addNewConfig(file *os.File, newConfig *appConfig) error {
	file.Seek(0, io.SeekStart)
	appsConfig := make(map[string]*appConfig)

	err := yaml.NewDecoder(file).Decode(appsConfig)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	// if no app configs already exist, make the new one the default
	if len(appsConfig) == 0 {
		newConfig.Default = true
	}

	if _, ok := appsConfig[newConfig.Name]; ok {
		return fmt.Errorf("configuration for %q already exists", newConfig.Name)
	}

	if newConfig.URL == "" {
		newConfig.URL = defaultEdgeURL
	}

	newCfg := map[string]*appConfig{
		newConfig.Name: newConfig,
	}
	err = yaml.NewEncoder(file).Encode(newCfg)
	return err
}

func removeConfigCmd() *cli.Command {
	cmd := &cli.Command{
		Name:        "remove",
		Usage:       "Remove an App configuration",
		UsageText:   "stream-cli config remove <app to remove>",
		Description: "Remove an App configuration. This operation is irrevocable",

		//SilenceUsage:  true,
		//SilenceErrors: true,

		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return fmt.Errorf("remove command accepts 1 argument")
			}

			f, err := getConfigurationFile()
			if err != nil {
				return err
			}
			defer f.Close()

			return removeConfig(f, ctx.Args().First())
		},
	}

	return cmd
}

func removeConfig(file *os.File, app string) error {
	appsConfig := make(map[string]*appConfig)

	err := yaml.NewDecoder(file).Decode(appsConfig)
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("config file is empty")
		}
		return err
	}

	var found bool
	filteredAppsConfig := make(map[string]*appConfig, len(appsConfig))
	for k, v := range appsConfig {
		if k == app {
			found = true
			continue
		}
		filteredAppsConfig[k] = v
	}

	if !found {
		return fmt.Errorf("application %q doesn't exist", app)
	}

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("cannot truncate configuration file")
	}
	file.Seek(0, io.SeekStart)

	if len(filteredAppsConfig) == 0 {
		return nil
	}
	err = yaml.NewEncoder(file).Encode(filteredAppsConfig)
	return err
}

func listConfigsCmd() *cli.Command {
	cmd := &cli.Command{
		Name:        "list",
		Usage:       "List all configurations",
		UsageText:   "stream-cli config list",
		Description: "List all app configurations",

		//SilenceUsage:  true,
		//SilenceErrors: true,

		Action: func(_ *cli.Context) error {
			f, err := getConfigurationFile()
			if err != nil {
				return err
			}
			defer f.Close()

			appsConfig := make(map[string]*appConfig)

			err = yaml.NewDecoder(f).Decode(appsConfig)
			if err != nil {
				if err == io.EOF {
					return fmt.Errorf("config file is empty")
				}
				return err
			}

			t := tabby.New()
			t.AddHeader("", "Name", "Access Key", "Secret Key", "Region")

			for k, v := range appsConfig {
				defaultApp := ""
				if v.Default {
					defaultApp = "(default)"
				}

				secret := fmt.Sprintf("**************%v", v.AccessSecretKey[len(v.AccessSecretKey)-4:])
				t.AddLine(defaultApp, k, v.AccessKey, secret, v.URL)
			}
			t.Print()
			return nil
		},
	}

	return cmd
}

func getConfigurationFile() (*os.File, error) {
	d, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot get user's home directory: %v", err)
	}

	err = os.Mkdir(path.Join(d, configDir), 0755)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("cannot create config directory: %v", err)
	}

	filepath := path.Join(d, configDir, configFile)
	return os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
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
			Transform: func(ans interface{}) interface{} {
				s, ok := ans.(string)
				if !ok {
					return defaultEdgeURL
				}
				if s == "" {
					return defaultEdgeURL
				}
				return s
			},
		},
	}
}
