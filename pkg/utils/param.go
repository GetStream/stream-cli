package utils

import (
	"encoding/json"
	"strings"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/spf13/pflag"
)

func GetJSONParam(f *pflag.FlagSet, name string) (map[string]any, error) {
	data, err := f.GetString(name)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	err := json.Unmarshal([]byte(data), &result)
	return result, err
}

func GetStringSliceParam(f *pflag.FlagSet, name string) ([]string, error) {
	data, err := f.GetString(name)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(data, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if v := strings.TrimSpace(part); v != "" {
			result = append(result, v)
		}
	}
	return result, nil
}

func GetPartialUpdateParam(f *pflag.FlagSet) (stream.PartialUpdate, error) {
	var update stream.PartialUpdate
	var err error
	update.Set, err = GetJSONParam(f, "set")
	if err != nil {
		return update, err
	}
	update.Unset, err = GetStringSliceParam(f, "unset")
	return update, err
}
