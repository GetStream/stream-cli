package version

import (
	"fmt"
)

const (
	versionMajor = 1
	versionMinor = 4
	versionPatch = 5
)

func FmtVersion() string {
	return fmt.Sprintf("%d.%d.%d",
		versionMajor,
		versionMinor,
		versionPatch)
}
