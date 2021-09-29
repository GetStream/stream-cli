package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddNewConfig(t *testing.T) {
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

	var buf bytes.Buffer
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := test.config()
			err := addNewConfig(&buf, &config)

			if test.errored {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expected, buf.String())
		})
	}
}
