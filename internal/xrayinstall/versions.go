// SPDX-License-Identifier: AGPL-3.0-or-later
package xrayinstall

import (
	"strconv"
	"strings"
)

// CompareVersions returns -1/0/1 like strings.Compare, on the SemVer
// MAJOR.MINOR.PATCH portion. Unknown / unparseable inputs sort last
// (returned as 0 vs known so the panel never claims a known version is
// older than an unknown one).
func CompareVersions(a, b string) int {
	pa, oka := parseSemver(a)
	pb, okb := parseSemver(b)
	if !oka && !okb {
		return 0
	}
	if !oka {
		return -1
	}
	if !okb {
		return 1
	}
	for k := 0; k < 3; k++ {
		if pa[k] < pb[k] {
			return -1
		}
		if pa[k] > pb[k] {
			return 1
		}
	}
	return 0
}

func parseSemver(s string) ([3]int, bool) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "v") {
		s = s[1:]
	}
	if i := strings.IndexAny(s, "-+"); i >= 0 {
		s = s[:i]
	}
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return [3]int{}, false
	}
	var out [3]int
	for k := 0; k < 3; k++ {
		if k >= len(parts) {
			out[k] = 0
			continue
		}
		n, err := strconv.Atoi(parts[k])
		if err != nil {
			return [3]int{}, false
		}
		out[k] = n
	}
	return out, true
}
