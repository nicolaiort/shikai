package version

import (
	"fmt"
	"strconv"
	"strings"
)

// Bump increments the version based on the bump type.
func Bump(current string, bump string) string {
	major, minor, patch := parseVersion(current)

	switch bump {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	}

	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

// AddPrerelease adds a prerelease suffix to the version.
func AddPrerelease(ver string, id string) string {
	return fmt.Sprintf("%s-%s.0", ver, id)
}

// parseVersion splits a semver version into major, minor, patch.
func parseVersion(ver string) (int, int, int) {
	parts := strings.Split(ver, ".")
	if len(parts) != 3 {
		return 0, 0, 0
	}
	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])
	return major, minor, patch
}
