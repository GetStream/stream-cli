package cli

import (
	"fmt"
)

const (
	versionMajor = 1
	versionMinor = 0
	versionPatch = 0
)

func fmtVersion() string {
	return fmt.Sprintf("%d.%d.%d",
		versionMajor,
		versionMinor,
		versionPatch)
}
