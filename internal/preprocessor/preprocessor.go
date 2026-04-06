package preprocessor

import (
	"regexp"
	"strings"
)

var (
	reSeparators = regexp.MustCompile(`(?m)^\s*(⸻|---|___|\*\*\*+)\s*$`)
	reGarbage    = regexp.MustCompile(`(?i)^\s*(explanation|solution|note|👉|✅).*`)
	reLevel      = regexp.MustCompile(`(?i)^\s*(📘.*|level[_\s]?\d+.*)$`)
)

func Process(input string) string {
	lines := strings.Split(input, "\n")
	var cleaned []string

	for _, line := range lines {
		t := strings.TrimSpace(line)

		if t == "" {
			continue
		}

		// ❌ remove garbage lines
		if reSeparators.MatchString(t) ||
			reGarbage.MatchString(t) ||
			reLevel.MatchString(t) {
			continue
		}

		cleaned = append(cleaned, line)
	}

	return strings.Join(cleaned, "\n")
}
