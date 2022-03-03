package config

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func getTempFile(t *testing.T) *os.File {
	tmpFile, err := os.CreateTemp("", "*.yaml")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Remove(tmpFile.Name())
	})
	return tmpFile
}

func TestAddNewConfig(t *testing.T) {
	tests := []struct {
		name      string
		appConfig func() App
		expected  string
		errored   bool
	}{
		{
			name: "add first configuration",
			appConfig: func() App {
				app := App{ChatURL: DefaultChatEdgeURL}
				app.Name = "BestConfig"
				app.AccessKey = "FamousKey"
				app.AccessSecretKey = "TopSecret"
				return app
			},
			expected: `
apps:
    - name: BestConfig
      access-key: FamousKey
      access-secret-key: TopSecret
      url: https://chat.stream-io-api.com
default: BestConfig
`,
		},
		{
			name: "add second configuration",
			appConfig: func() App {
				app := App{ChatURL: DefaultChatEdgeURL}
				app.Name = "BestConfigEver"
				app.AccessKey = "FamousKey"
				app.AccessSecretKey = "TopSecret"
				return app
			},
			expected: `
apps:
    - name: BestConfig
      access-key: FamousKey
      access-secret-key: TopSecret
      url: https://chat.stream-io-api.com
    - name: BestConfigEver
      access-key: FamousKey
      access-secret-key: TopSecret
      url: https://chat.stream-io-api.com
default: BestConfig
`,
		},
		{
			name: "add already existing configuration",
			appConfig: func() App {
				app := App{ChatURL: DefaultChatEdgeURL}
				app.Name = "BestConfig"
				app.AccessKey = "FamousKey"
				app.AccessSecretKey = "TopSecret"
				return app
			},
			errored: true,
		},
	}

	file := getTempFile(t)
	viper.SetConfigFile(file.Name())
	config := &Config{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newApp := test.appConfig()
			err := config.Add(newApp)
			if test.errored {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			content, err := os.ReadFile(file.Name())
			require.NoError(t, err)

			require.Equal(t, getNormalizedString(test.expected), getNormalizedString(string(content)))
		})
	}
}

func TestRemoveConfig(t *testing.T) {
	file := getTempFile(t)
	viper.SetConfigFile(file.Name())
	config := &Config{}

	err := config.Add(App{
		Name:            "test1",
		AccessKey:       "test1",
		AccessSecretKey: "test1",
		ChatURL:         DefaultChatEdgeURL,
	})
	require.NoError(t, err)

	err = config.Add(App{
		Name:            "test2",
		AccessKey:       "test2",
		AccessSecretKey: "test2",
		ChatURL:         DefaultChatEdgeURL,
	})
	require.NoError(t, err)

	// remove non-existing app configuration should fail
	err = config.Remove("unknown")
	require.Error(t, err)

	err = config.Remove("test1")
	require.NoError(t, err)

	expected := `
apps:
    - name: test2
      access-key: test2
      access-secret-key: test2
      url: https://chat.stream-io-api.com
default: ""
`
	content, err := os.ReadFile(file.Name())
	require.NoError(t, err)
	require.Equal(t, getNormalizedString(expected), getNormalizedString(string(content)))
}

func TestSetDefault(t *testing.T) {
	file := getTempFile(t)
	viper.SetConfigFile(file.Name())
	config := &Config{}

	err := config.Add(App{
		Name:            "test1",
		AccessKey:       "test1",
		AccessSecretKey: "test1",
		ChatURL:         DefaultChatEdgeURL,
	})
	require.NoError(t, err)

	err = config.Add(App{
		Name:            "test2",
		AccessKey:       "test2",
		AccessSecretKey: "test2",
		ChatURL:         DefaultChatEdgeURL,
	})
	require.NoError(t, err)

	require.True(t, config.Default == "test1")

	err = config.SetDefault("test2")
	require.NoError(t, err)

	require.True(t, config.Default == "test2")

	expected := `
apps:
    - name: test1
      access-key: test1
      access-secret-key: test1
      url: https://chat.stream-io-api.com
    - name: test2
      access-key: test2
      access-secret-key: test2
      url: https://chat.stream-io-api.com
default: test2
`

	content, err := os.ReadFile(file.Name())
	require.NoError(t, err)
	require.Equal(t, getNormalizedString(expected), getNormalizedString(string(content)))
}

func getNormalizedString(s string) string {
	noSpace := strings.Replace(s, " ", "", -1)
	noNewLine := strings.Replace(noSpace, "\n", "", -1)

	return strings.TrimSpace(noNewLine)
}
