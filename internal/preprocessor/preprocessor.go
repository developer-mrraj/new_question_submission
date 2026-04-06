package preprocessor

import (
	"regexp"
	"strings"
)

var (
	reSeparators = regexp.MustCompile(`(?m)^\s*(⸻|---|___|[\*]{3,})\s*$`)
	// Strips standalone explanation/solution/note lines from the MCQ (Other) flow
	reGarbage = regexp.MustCompile(`(?i)^\s*(?:explanation|solution|note)`)
	reLevel   = regexp.MustCompile(`(?i)^\s*(📘.*|level[_\s]?\d+.*)$`)
	// Detects the answer keyword so ✅ Answer:/👉 Final Answer: lines are preserved
	reAnswerKeyword = regexp.MustCompile(`(?i)(?:ans|answer)`)

	// reGarbageNonExplanation strips solution/note lines but keeps "Explanation:" lines.
	// Used by the DQ (Decision & Quiz) flow where Explanation is a meaningful field.
	reGarbageNonExplanation = regexp.MustCompile(`(?i)^\s*(?:solution|note)`)
)

// Process cleans raw MCQ text for the standard "Other" flow.
// It strips Explanation, separator, and junk lines.
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

		// ✅ emoji lines: keep ONLY if they contain an answer keyword
		if strings.HasPrefix(t, "✅") || strings.HasPrefix(t, "👉") {
			if !reAnswerKeyword.MatchString(t) {
				continue
			}
		}

		cleaned = append(cleaned, line)
	}

	return strings.Join(cleaned, "\n")
}

// ProcessDQ cleans raw text for the Decision & Quiz flow.
// Same as Process but PRESERVES "Explanation: ..." lines so the DQ parser can extract them.
func ProcessDQ(input string) string {
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

		// ❌ Remove solution/note lines — but NOT explanation lines
		if reGarbageNonExplanation.MatchString(t) {
			continue
		}

		// ❌ Remove level headers
		if reLevel.MatchString(t) {
			continue
		}

		// ✅ emoji lines: keep ONLY if they contain an answer keyword
		if strings.HasPrefix(t, "✅") || strings.HasPrefix(t, "👉") {
			if !reAnswerKeyword.MatchString(t) {
				continue
			}
		}

		cleaned = append(cleaned, line)
	}

	return strings.Join(cleaned, "\n")
}
