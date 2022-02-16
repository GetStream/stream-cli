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
	"gopkg.in/yaml.v3"
)

const (
	configDir  = "stream-cli"
	configFile = "config.yml"

	defaultEdgeURL = "https://chat.stream-io-api.com"
)

func NewRootConfigCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "config",
		Usage:       "Manage app configurations",
		Description: `Manage app configurations`,
		Subcommands: []*cli.Command{
			newConfigCmd(config),
			removeConfigCmd(config),
			listConfigsCmd(config),
			defaultConfigCmd(config),
		},
	}
}

func newConfigCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "new",
		Usage:       "Add a new App configuration",
		UsageText:   "stream-cli config new",
		Description: "Add a new App configuration which can be used on further operations",

		Action: func(ctx *cli.Context) error {
			newConfig := newDefaultConfig()
			err := survey.Ask(questions(), &newConfig)
			if err != nil {
				return err
			}

			return config.Add(newConfig)
		},
	}
}

func removeConfigCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "remove",
		Usage:       "Remove an App configuration",
		UsageText:   "stream-cli config remove <app to remove>",
		Description: "Remove an App configuration. This operation is irrevocable",

		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return errors.New("remove command requires 1 argument")
			}
			return config.Remove(ctx.Args().First())
		},
	}
}

func listConfigsCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "List all configurations",
		UsageText:   "stream-cli config list",
		Description: "List all app configurations",

		Action: func(_ *cli.Context) error {
			t := tabby.New()
			t.AddHeader("", "Name", "Access Key", "Secret Key", "Region")

			for k, v := range config.appsConfig {
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
}

func defaultConfigCmd(config *Config) *cli.Command {
	cmd := &cli.Command{
		Name:        "default",
		Usage:       "Set a configuration as the default",
		UsageText:   "stream-cli config default <name of the configuration>",
		Description: "Set a configuration as the default, it will be used by default for all operations",

		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return errors.New("default command requires 1 argument")
			}
			return config.SetDefault(ctx.Args().First())
		},
	}

	return cmd
}

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

type Config struct {
	appsConfig map[string]*appConfig
	file       *os.File
}

func NewConfig() (*Config, error) {
	d, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("cannot get user's home directory: %v", err)
	}

	err = os.Mkdir(path.Join(d, configDir), 0755)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("cannot create config directory: %v", err)
	}

	filepath := path.Join(d, configDir, configFile)
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	appsConfig := make(map[string]*appConfig)
	err = yaml.NewDecoder(file).Decode(appsConfig)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	return &Config{
		file: file,

		appsConfig: appsConfig,
	}, nil
}

func (c *Config) Add(newConfig appConfig) error {
	if len(c.appsConfig) == 0 {
		newConfig.Default = true
	}

	if _, ok := c.appsConfig[newConfig.Name]; ok {
		return fmt.Errorf("configuration for %q already exists", newConfig.Name)
	}

	if newConfig.URL == "" {
		newConfig.URL = defaultEdgeURL
	}

	if c.appsConfig == nil {
		c.appsConfig = make(map[string]*appConfig)
	}
	c.appsConfig[newConfig.Name] = &newConfig
	return c.Append(newConfig)
}

func (c *Config) Remove(configName string) error {
	if len(c.appsConfig) == 0 {
		return errors.New("config file is empty")
	}

	if _, ok := c.appsConfig[configName]; !ok {
		return fmt.Errorf("application %q doesn't exist", configName)
	}

	delete(c.appsConfig, configName)

	return c.Override()
}

func (c *Config) SetDefault(configName string) error {
	if len(c.appsConfig) == 0 {
		return errors.New("config file is empty")
	}

	config, ok := c.appsConfig[configName]
	if !ok {
		return fmt.Errorf("application %q doesn't exist", configName)
	}

	if config.Default {
		// if already default, early return
		return nil
	}

	for k, v := range c.appsConfig {
		if k == configName {
			v.Default = true
			continue
		}
		v.Default = false
	}

	return c.Override()
}

func (c *Config) Append(newCfg appConfig) error {
	newAppConfig := map[string]*appConfig{
		newCfg.Name: &newCfg,
	}
	return yaml.NewEncoder(c.file).Encode(newAppConfig)
}

func (c *Config) Override() error {
	err := c.file.Truncate(0)
	if err != nil {
		return errors.New("cannot truncate configuration file")
	}
	c.file.Seek(0, io.SeekStart)

	if len(c.appsConfig) == 0 {
		return nil
	}
	return yaml.NewEncoder(c.file).Encode(c.appsConfig)
}

func (c *Config) Close() {
	if c.file != nil {
		c.file.Close()
	}
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
