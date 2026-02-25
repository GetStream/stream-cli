package version

import (
	"fmt"
)

const (
	versionMajor = 1
	versionMinor = 8
	versionPatch = 4
)

func FmtVersion() string {
	return fmt.Sprintf("%d.%d.%d",
		versionMajor,
		versionMinor,
		versionPatch)
}
