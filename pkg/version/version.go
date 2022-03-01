package version

import (
	"fmt"
)

const (
	versionMajor = 1
	versionMinor = 0
	versionPatch = 0
)

func FmtVersion() string {
	return fmt.Sprintf("%d.%d.%d",
		versionMajor,
		versionMinor,
		versionPatch)
}
