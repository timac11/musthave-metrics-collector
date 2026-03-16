package version

import (
	"fmt"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func format(val string) string {
	if val == "" {
		return "N/A"
	}

	return val
}

func Print() {
	fmt.Printf("Build version: %s\n", format(buildVersion))
	fmt.Printf("Build date: %s\n", format(buildDate))
	fmt.Printf("Build commit: %s\n", format(buildCommit))
}
