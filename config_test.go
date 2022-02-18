package cli

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func getFile(t *testing.T) *os.File {
	tmpFile, err := os.CreateTemp("", "testconfig.yaml")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})
	return tmpFile
}

func TestAddNewConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   func() appConfig
		expected string
		errored  bool
	}{
		{
			name: "add first configuration",
			config: func() appConfig {
				cfg := newDefaultConfig()
				cfg.Name = "BestConfig"
				cfg.AccessKey = "FamousKey"
				cfg.AccessSecretKey = "TopSecret"
				return cfg
			},
			expected: `BestConfig:
    access-key: FamousKey
    access-secret-key: TopSecret
    url: https://chat.stream-io-api.com
    default: true
`,
		},
		{
			name: "add second configuration",
			config: func() appConfig {
				cfg := newDefaultConfig()
				cfg.Name = "BestConfigEver"
				cfg.AccessKey = "FamousKey"
				cfg.AccessSecretKey = "TopSecret"
				return cfg
			},
			expected: `BestConfig:
    access-key: FamousKey
    access-secret-key: TopSecret
    url: https://chat.stream-io-api.com
    default: true
BestConfigEver:
    access-key: FamousKey
    access-secret-key: TopSecret
    url: https://chat.stream-io-api.com
`,
		},
		{
			name: "add already existing configuration",
			config: func() appConfig {
				cfg := newDefaultConfig()
				cfg.Name = "BestConfig"
				cfg.AccessKey = "FamousKey"
				cfg.AccessSecretKey = "TopSecret"
				return cfg
			},
			errored: true,
		},
	}

	file := getFile(t)
	config := &Config{
		filepath: file.Name(),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conf := test.config()
			err := config.Add(conf)
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
		filepath: file.Name(),
	}

	err := config.Add(appConfig{
		Name:            "test1",
		AccessKey:       "test1",
		AccessSecretKey: "test1",
		URL:             defaultEdgeURL,
	})
	require.NoError(t, err)

	err = config.Add(appConfig{
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

	expected := `test2:
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
		filepath: file.Name(),
	}

	err := config.Add(appConfig{
		Name:            "test1",
		AccessKey:       "test1",
		AccessSecretKey: "test1",
		URL:             defaultEdgeURL,
	})
	require.NoError(t, err)

	err = config.Add(appConfig{
		Name:            "test2",
		AccessKey:       "test2",
		AccessSecretKey: "test2",
		URL:             defaultEdgeURL,
	})
	require.NoError(t, err)

	require.True(t, config.appsConfig["test1"].Default)
	require.False(t, config.appsConfig["test2"].Default)

	err = config.SetDefault("test2")
	require.NoError(t, err)

	require.False(t, config.appsConfig["test1"].Default)
	require.True(t, config.appsConfig["test2"].Default)

	expected := `test1:
    access-key: test1
    access-secret-key: test1
    url: https://chat.stream-io-api.com
test2:
    access-key: test2
    access-secret-key: test2
    url: https://chat.stream-io-api.com
    default: true
`
	content, err := os.ReadFile(file.Name())
	require.NoError(t, err)
	require.Equal(t, expected, string(content))
}
