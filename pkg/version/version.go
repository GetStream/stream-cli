package version

import (
	"fmt"
)

const (
	versionMajor = 0
	versionMinor = 0
	versionPatch = 1
)

func FmtVersion() string {
	return fmt.Sprintf("%d.%d.%d",
		versionMajor,
		versionMinor,
		versionPatch)
}
