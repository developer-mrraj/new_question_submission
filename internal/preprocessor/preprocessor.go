package preprocessor

import (
	"regexp"
	"strings"
)

var (
	reSeparators = regexp.MustCompile(`(?m)^\s*(⸻|---|___|\*\*\*+)\s*$`)
	// Strips standalone hint/explanation lines (NOT ✅ or 👉 — handled separately below)
	reGarbage = regexp.MustCompile(`(?i)^\s*(?:explanation|solution|note)`)
	reLevel   = regexp.MustCompile(`(?i)^\s*(📘.*|level[_\s]?\d+.*)$`)
	// Detects the answer keyword so ✅ Answer:/👉 Final Answer: lines are preserved
	reAnswerKeyword = regexp.MustCompile(`(?i)(?:ans|answer)`)
)

func Process(input string) string {
	lines := strings.Split(input, "\n")
	var cleaned []string

	for _, line := range lines {
		t := strings.TrimSpace(line)

		if t == "" {
			continue
		}

		// ❌ Remove separator lines
		if reSeparators.MatchString(t) {
			continue
		}

		// ❌ Remove explanation/hint lines
		if reGarbage.MatchString(t) {
			continue
		}

		// ❌ Remove level headers
		if reLevel.MatchString(t) {
			continue
		}

		// ✅ emoji lines: keep ONLY if they contain an answer keyword (e.g. "✅ Ans: B)", "👉 Answer: C)")
		// otherwise strip (they are decorative lines like "✅ Key: ...")
		if strings.HasPrefix(t, "✅") || strings.HasPrefix(t, "👉") {
			if !reAnswerKeyword.MatchString(t) {
				continue
			}
		}

		cleaned = append(cleaned, line)
	}

	return strings.Join(cleaned, "\n")
}
