package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v2"
	"golang.org/x/mod/semver"
)

func NewUpdateCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "update",
		Usage:       "Handles self-updates of the Stream CLI.",
		Description: "Handles self-updates of the Stream CLI.",
		Subcommands: []*cli.Command{
			selfUpdateCmd(config),
		},
	}
}

func selfUpdateCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "self",
		Usage:       "Self-updates of the Stream CLI if there is a newer version available.",
		UsageText:   "stream-cli update self",
		Description: "Self-updates of the Stream CLI if there is a newer version available.",
		Action: func(c *cli.Context) error {
			current := fmtVersion()
			latest, err := latestVersion()
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if semver.Compare(latest, current) == 1 {
				if err = updateBinary(latest); err != nil {
					return cli.Exit(err.Error(), 1)
				}
				PrintHappyMessageFormatted("Successfully updated Stream CLI to v%s", latest)
			} else {
				PrintHappyMessage("You are already on the latest Stream CLI version.")
			}

			return nil
		},
	}
}

func latestVersion() (string, error) {
	resp, err := http.DefaultClient.Get("https://api.github.com/repos/getstream/stream-cli/releases/latest")

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	releases := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", releases["tag_name"])[1:], nil
}

func updateBinary(latest string) error {
	var u string

	switch runtime.GOOS {
	case "windows":
		u = "https://github.com/GetStream/stream-cli/releases/latest/download/windows.exe"
	case "darwin":
		u = "https://github.com/GetStream/stream-cli/releases/latest/download/darwin"
	case "linux":
		u = "https://github.com/GetStream/stream-cli/releases/latest/download/linux"
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	resp, err := http.DefaultClient.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return overWriteBinary(resp)
}

func overWriteBinary(resp *http.Response) error {
	path, err := os.Executable()
	if err != nil {
		return err
	}
	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(realPath, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
