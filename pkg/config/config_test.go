package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func getFile(t *testing.T) *os.File {
	tmpFile, err := os.CreateTemp("", "testconfig.yaml")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Remove(tmpFile.Name())
	})
	return tmpFile
}

func TestNewConfig(t *testing.T) {
	c, err := NewConfig(os.TempDir())
	require.NoError(t, err)
	require.True(t, strings.HasSuffix(c.FilePath, filepath.Join(configDir, configFile)))
	_ = os.RemoveAll(filepath.Join(os.TempDir(), configDir))
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
				app := App{URL: defaultEdgeURL}
				app.Name = "BestConfig"
				app.AccessKey = "FamousKey"
				app.AccessSecretKey = "TopSecret"
				return app
			},
			expected: `default: BestConfig
apps:
    - name: BestConfig
      access-key: FamousKey
      access-secret-key: TopSecret
      url: https://chat.stream-io-api.com
`,
		},
		{
			name: "add second configuration",
			appConfig: func() App {
				app := App{URL: defaultEdgeURL}
				app.Name = "BestConfigEver"
				app.AccessKey = "FamousKey"
				app.AccessSecretKey = "TopSecret"
				return app
			},
			expected: `default: BestConfig
apps:
    - name: BestConfig
      access-key: FamousKey
      access-secret-key: TopSecret
      url: https://chat.stream-io-api.com
    - name: BestConfigEver
      access-key: FamousKey
      access-secret-key: TopSecret
      url: https://chat.stream-io-api.com
`,
		},
		{
			name: "add already existing configuration",
			appConfig: func() App {
				app := App{URL: defaultEdgeURL}
				app.Name = "BestConfig"
				app.AccessKey = "FamousKey"
				app.AccessSecretKey = "TopSecret"
				return app
			},
			errored: true,
		},
	}

	file := getFile(t)
	config := &Config{
		FilePath: file.Name(),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newApp := test.appConfig()
			err := config.Add(newApp)
			if test.errored {
				require.Error(t, err)
				return
			}

			content, err := os.ReadFile(file.Name())
			require.NoError(t, err)

			require.NoError(t, err)
			require.Equal(t, test.expected, string(content))
		})
	}
}

func TestRemoveConfig(t *testing.T) {
	file := getFile(t)
	config := &Config{
		FilePath: file.Name(),
	}

	err := config.Add(App{
		Name:            "test1",
		AccessKey:       "test1",
		AccessSecretKey: "test1",
		URL:             defaultEdgeURL,
	})
	require.NoError(t, err)

	err = config.Add(App{
		Name:            "test2",
		AccessKey:       "test2",
		AccessSecretKey: "test2",
		URL:             defaultEdgeURL,
	})
	require.NoError(t, err)

	// remove non-existing app configuration should fail
	err = config.Remove("unknown")
	require.Error(t, err)

	err = config.Remove("test1")
	require.NoError(t, err)

	expected := `default: ""
apps:
    - name: test2
      access-key: test2
      access-secret-key: test2
      url: https://chat.stream-io-api.com
`
	content, err := os.ReadFile(file.Name())
	require.NoError(t, err)
	require.Equal(t, expected, string(content))
}

func TestSetDefault(t *testing.T) {
	file := getFile(t)
	config := &Config{
		FilePath: file.Name(),
	}

	err := config.Add(App{
		Name:            "test1",
		AccessKey:       "test1",
		AccessSecretKey: "test1",
		URL:             defaultEdgeURL,
	})
	require.NoError(t, err)

	err = config.Add(App{
		Name:            "test2",
		AccessKey:       "test2",
		AccessSecretKey: "test2",
		URL:             defaultEdgeURL,
	})
	require.NoError(t, err)

	require.True(t, config.Default == "test1")

	err = config.SetDefault("test2")
	require.NoError(t, err)

	require.True(t, config.Default == "test2")

	expected := `default: test2
apps:
    - name: test1
      access-key: test1
      access-secret-key: test1
      url: https://chat.stream-io-api.com
    - name: test2
      access-key: test2
      access-secret-key: test2
      url: https://chat.stream-io-api.com
`

	content, err := os.ReadFile(file.Name())
	require.NoError(t, err)
	require.Equal(t, expected, string(content))
}
