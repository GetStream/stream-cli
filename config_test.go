package cli

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddNewConfig(t *testing.T) {
	tmpFile, err := os.Create("testconfig.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	tests := []struct{
		name string
		config func() appConfig
		expected string
		errored bool
	}{
		{
			name: "add first configuration",
			config: func()appConfig{
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
`,
		},
		{
			name: "add second configuration",
			config: func()appConfig{
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
BestConfigEver:
  access-key: FamousKey
  access-secret-key: TopSecret
  url: https://chat.stream-io-api.com
`,
		},
		{
			name: "add already existing configuration",
			config: func()appConfig{
				cfg := newDefaultConfig()
				cfg.Name = "BestConfig"
				cfg.AccessKey = "FamousKey"
				cfg.AccessSecretKey = "TopSecret"
				return cfg
			},
			errored: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := test.config()
			err := addNewConfig(tmpFile, &config)

			if test.errored {
				require.Error(t, err)
				return
			}

			content, err := ioutil.ReadFile(tmpFile.Name())
			require.NoError(t, err)

			require.NoError(t, err)
			require.Equal(t, test.expected, string(content))
		})
	}
}

func TestRemoveConfig(t *testing.T) {
	tmpFile, err := os.OpenFile("testconfig.yaml", os.O_RDWR|os.O_CREATE, 0644)
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	config := appConfig{
		Name: "test1",
		AccessKey: "test1",
		AccessSecretKey: "test1",
		URL: defaultEdgeURL,
	}
	err = addNewConfig(tmpFile, &config)
	require.NoError(t, err)

	config.Name = "test2"
	config.AccessKey = "test2"
	config.AccessSecretKey = "test2"
	config.URL = defaultEdgeURL
	err = addNewConfig(tmpFile, &config)
	require.NoError(t, err)

	// remove non-existing app configuration should fail
	err = removeConfig(tmpFile, "unknown")
	require.Error(t, err)

	err = removeConfig(tmpFile, "test1")
	require.NoError(t, err)

	expected := `test2:
  access-key: test2
  access-secret-key: test2
  url: https://chat.stream-io-api.com
`
	content, err := ioutil.ReadFile(tmpFile.Name())
	require.NoError(t, err)
	require.Equal(t, expected, string(content))
}
