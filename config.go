package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	stream "github.com/GetStream/stream-chat-go/v5"
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
		Usage:       "Manage app configurations.",
		Description: "Manage app configurations.",
		Subcommands: []*cli.Command{
			newAppCmd(config),
			removeAppCmd(config),
			listAppsCmd(config),
			setAppDefaultCmd(config),
		},
	}
}

func newAppCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "new",
		Usage:       "Add a new application",
		UsageText:   "stream-cli config new",
		Description: "Add a new application which can be used for further operations",

		Action: func(ctx *cli.Context) error {
			return RunQuestionnaire(ctx, config)
		},
	}
}

func RunQuestionnaire(ctx *cli.Context, config *Config) error {
	var newAppConfig App
	err := survey.Ask(questions(), &newAppConfig)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	err = config.Add(newAppConfig)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	PrintMessage(ctx, "Application successfully added. 🚀")
	return nil
}

func removeAppCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "remove",
		Usage:       "Remove an application",
		UsageText:   "stream-cli config remove <app to remove>",
		Description: "Remove an application. This operation is irrevocable",

		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return errors.New("remove command requires 1 argument")
			}
			return config.Remove(ctx.Args().First())
		},
	}
}

func listAppsCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "List all applications",
		UsageText:   "stream-cli config list",
		Description: "List all applications configurations",

		Action: func(_ *cli.Context) error {
			t := tabby.New()
			t.AddHeader("", "Name", "Access Key", "Secret Key", "Region")

			defaultApp := config.Default
			for _, app := range config.Apps {
				def := ""
				if app.Name == defaultApp {
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

func setAppDefaultCmd(config *Config) *cli.Command {
	cmd := &cli.Command{
		Name:        "default",
		Usage:       "Set an application as the default",
		UsageText:   "stream-cli config default <name of the application>",
		Description: "Set an application as the default, it will be used by default for all operations",

		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() != 1 {
				return errors.New("default command requires 1 argument")
			}
			return config.SetDefault(ctx.Args().First())
		},
	}

	return cmd
}

type Config struct {
	// Default is the default configuration used for operations
	Default string `yaml:"default"`
	Apps    []App  `yaml:"apps"`

	filePath string
}

type App struct {
	Name            string `yaml:"name"`
	AccessKey       string `yaml:"access-key"`
	AccessSecretKey string `yaml:"access-secret-key"`
	URL             string `yaml:"url"`
}

func NewConfig(dir string) (*Config, error) {
	err := os.Mkdir(filepath.Join(dir, configDir), 0755)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("cannot create configuration directory: %v", err)
	}

	fp := filepath.Join(dir, configDir, configFile)
	b, err := os.ReadFile(fp)
	// don't return an error if the file doesn't exist yet, it may not exist on the first run
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	config := &Config{
		filePath: fp,
	}
	err = yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) Get(name string) (*App, error) {
	if len(c.Apps) == 0 {
		return nil, errors.New("no application configured, please run `stream-cli config new` to add a new one")
	}

	for _, app := range c.Apps {
		if app.Name == name {
			return &app, nil
		}
	}
	return nil, fmt.Errorf("application %q doesn't exist", name)
}

func (c *Config) GetCredentials(ctx *cli.Context) (string, string, error) {
	appName := c.Default
	if explicit := ctx.String("app"); explicit != "" {
		appName = explicit
	}

	a, err := c.Get(appName)
	if err != nil {
		return "", "", err
	}

	return a.AccessKey, a.AccessSecretKey, nil
}

func (c *Config) GetStreamClient(ctx *cli.Context) (*stream.Client, error) {
	key, secret, err := c.GetCredentials(ctx)
	if err != nil {
		return nil, err
	}

	return stream.NewClient(key, secret)
}

func (c *Config) Add(newApp App) error {
	if len(c.Apps) == 0 {
		c.Default = newApp.Name
	}

	for _, app := range c.Apps {
		if newApp.Name == app.Name {
			return fmt.Errorf("application %q already exists", newApp.Name)
		}
	}

	if newApp.URL == "" {
		newApp.URL = defaultEdgeURL
	}

	c.Apps = append(c.Apps, newApp)
	return c.WriteToFile()
}

func (c *Config) Remove(appName string) error {
	var (
		idx   int
		found bool
	)
	for i, app := range c.Apps {
		if appName == app.Name {
			found = true
			idx = i
			break
		}
	}
	if !found {
		return fmt.Errorf("application %q doesn't exist", appName)
	}

	if c.Default == appName {
		c.Default = ""
	}

	c.Apps = append(c.Apps[:idx], c.Apps[idx+1:]...)
	return c.WriteToFile()
}

func (c *Config) SetDefault(appName string) error {
	if c.Default == appName {
		// if already default, early return
		return nil
	}

	var found bool
	for _, app := range c.Apps {
		if appName == app.Name {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("application %q doesn't exist", appName)
	}

	c.Default = appName
	return c.WriteToFile()
}

func (c *Config) WriteToFile() error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(c.filePath, b, 0644)
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

				if !ok || s == "" {
					return defaultEdgeURL
				}

				return s
			},
		},
	}
}
