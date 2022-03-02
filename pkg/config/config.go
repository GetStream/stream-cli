package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configDir  = "stream-cli"
	configFile = "config.yml"

	DefaultEdgeURL = "https://chat.stream-io-api.com"
)

type Config struct {
	// It's so stupid but Viper uses `yaml` for deserialization into a map object
	// and then uses `mapstructure` to deserialize into a an actual Config object.

	// Default is the default configuration used for operations
	Default string `yaml:"default" mapstructure:"default"`
	Apps    []App  `yaml:"apps" mapstructure:"apps"`
}

type App struct {
	Name            string `yaml:"name" mapstructure:"name"`
	AccessKey       string `yaml:"access-key" mapstructure:"access-key"`
	AccessSecretKey string `yaml:"access-secret-key" mapstructure:"access-secret-key"`
	URL             string `yaml:"url" mapstructure:"url"`
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

func (c *Config) getCredentials(cmd *cobra.Command) (string, string, error) {
	appName := c.Default
	explicit, err := cmd.Flags().GetString("app")
	if err != nil {
		return "", "", err
	}
	if explicit != "" {
		appName = explicit
	}

	a, err := c.Get(appName)
	if err != nil {
		return "", "", err
	}

	return a.AccessKey, a.AccessSecretKey, nil
}

func (c *Config) GetStreamClient(cmd *cobra.Command) (*stream.Client, error) {
	key, secret, err := c.getCredentials(cmd)
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
		newApp.URL = DefaultEdgeURL
	}

	c.Apps = append(c.Apps, newApp)

	viper.Set("default", c.Default)
	viper.Set("apps", c.Apps)
	return viper.WriteConfig()
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

	viper.Set("default", c.Default)
	viper.Set("apps", c.Apps)
	return viper.WriteConfig()
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
	viper.Set("default", c.Default)
	return viper.WriteConfig()
}

func GetConfig(cmd *cobra.Command) *Config {
	cfg := &Config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		cmd.PrintErr("Configuration is malformed: " + err.Error())
		os.Exit(1)
	}

	return cfg
}

func GetInitConfig(cmd *cobra.Command, cfgPath *string) func() {
	return func() {
		var configPath string

		if *cfgPath != "" {
			// Use config file from the flag.
			configPath = *cfgPath
		} else {
			// Otherwise use UserConfigDir
			dir, err := os.UserConfigDir()
			if err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}

			configPath = filepath.Join(dir, configDir, "config.yml")
		}

		viper.SetConfigFile(configPath)
		viper.SetConfigType("yaml")

		err := viper.ReadInConfig()
		if err != nil && os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(configPath), 0755)
			if err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}

			f, err := os.Create(configPath)
			if err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}

			f.Close()
		}
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
	}

}
