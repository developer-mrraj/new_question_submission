package splitter

import (
	"regexp"
	"strings"
)

type QuestionBlock struct {
	Number int
	Raw    string
}

// ✅ Detect ALL formats:
// ✅ Detect ALL formats:
// 1.
// Q1.
// q1.
// 1)
// Q1)
// 🔥 Q19
var reQStart = regexp.MustCompile(`(?i)^\s*(?:🔥\s*)?(?:q\s*\d+(?:\s*[\.\)])?|\d+\s*[\.\)])`)

// REMOVE level / garbage
var reJunk = regexp.MustCompile(`(?i)^\s*(level|📘|⸻|---|___)`)

func Split(text string) []QuestionBlock {
	lines := strings.Split(text, "\n")

	var blocks []QuestionBlock
	var current []string
	qNum := 0

	flush := func() {
		if len(current) == 0 {
			return
		}

		raw := strings.TrimSpace(strings.Join(current, "\n"))

		if raw != "" {
			qNum++
			blocks = append(blocks, QuestionBlock{
				Number: qNum,
				Raw:    raw,
			})
		}

		current = nil
	}

	for _, line := range lines {
		t := strings.TrimSpace(line)

		// ❌ skip junk
		if reJunk.MatchString(t) {
			continue
		}

		// ✅ new question
		if reQStart.MatchString(t) {
			flush()
			current = []string{t}
			continue
		}

		// append
		if len(current) > 0 {
			current = append(current, t)
		}
	}

	flush()
	return blocks
}
