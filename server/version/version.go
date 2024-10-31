package version

import (
	"fmt"
	"golang.org/x/mod/semver"
)

func IsVersionGreaterOrEqualThan(version, target string) bool {
	return semver.Compare(fmt.Sprintf("v%s", version), fmt.Sprintf("v%s", target)) > -1
}

type SortVersion []string

func (s SortVersion) Len() int {
	return len(s)
}

func (s SortVersion) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortVersion) Less(i, j int) bool {
	v1 := fmt.Sprintf("v%s", s[i])
	v2 := fmt.Sprintf("v%s", s[j])
	return semver.Compare(v1, v2) == -1
}
