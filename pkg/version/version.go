package version

import (
	"fmt"
)

const (
	versionMajor = 1
	versionMinor = 7
	versionPatch = 1
)

func FmtVersion() string {
	return fmt.Sprintf("%d.%d.%d",
		versionMajor,
		versionMinor,
		versionPatch)
}
