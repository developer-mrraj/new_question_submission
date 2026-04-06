package parser

import (
	"regexp"
	"strings"
)

// ParsedDQQuestion is the same as ParsedQuestion but includes Title and Explanation.
// Used by the Decision & Quiz convert flow.
type ParsedDQQuestion struct {
	Number      int               `json:"number"`
	Title       string            `json:"title"`
	Text        string            `json:"text"`
	Explanation string            `json:"explanation"`
	Options     map[string]string `json:"options"`
	Answer      string            `json:"answer"`
}

// reExplanation matches lines like "Explanation: ..." (case-insensitive)
var reExplanation = regexp.MustCompile(`(?i)^explanation:\s*(.+)`)

// reTitleBracket extracts a topic inside parentheses at the START of the first line,
// e.g. "(Cleanliness & Responsibility)" → "Cleanliness & Responsibility"
// After the Q-header is stripped, the first line is like "(Cleanliness & Responsibility)"
// or "(Cleanliness & Responsibility) rest of question…"
var reTitleBracket = regexp.MustCompile(`^\(([^)]+)\)\s*(.*)`)

// ParseDQ parses a single question block and also extracts:
//   - Title: the topic in parentheses on the first line, e.g. "(Cleanliness & Responsibility)"
//   - Explanation: the "Explanation: ..." line if present
func ParseDQ(num int, raw string) ParsedDQQuestion {
	lines := strings.Split(raw, "\n")

	pq := ParsedDQQuestion{
		Number:  num,
		Options: map[string]string{},
	}

	var qLines []string
	var explanationLines []string
	optionsStarted := false

	for i, line := range lines {
		t := strings.TrimSpace(line)

		if t == "" {
			continue
		}

		// ── Line 0: strip Q-number prefix, then extract (Title) ─────────────
		if i == 0 {
			t = reQHeader.ReplaceAllString(t, "")
			t = strings.TrimSpace(t)

			// Try to extract "(Title) remaining question text" from line 0
			if m := reTitleBracket.FindStringSubmatch(t); m != nil {
				pq.Title = strings.TrimSpace(m[1]) // e.g. "Cleanliness & Responsibility"
				rest := strings.TrimSpace(m[2])
				if rest != "" {
					// Rare case: question text starts on the same line after the bracket
					qLines = append(qLines, rest)
				}
			} else {
				// No bracket — entire remaining text is part of the question
				if t != "" {
					qLines = append(qLines, t)
				}
			}
			continue
		}

		// ── Explanation line ─────────────────────────────────────────────────
		if m := reExplanation.FindStringSubmatch(t); m != nil {
			explanationLines = append(explanationLines, strings.TrimSpace(m[1]))
			continue
		}

		// ── Answer line ──────────────────────────────────────────────────────
		if m := reAnswer.FindStringSubmatch(t); m != nil {
			pq.Answer = strings.ToUpper(m[1])
			continue
		}

		// ── Option line (A., B), A-, etc.) ───────────────────────────────────
		if m := reOption.FindStringSubmatch(t); m != nil {
			opt := strings.ToUpper(m[1])
			val := strings.TrimSpace(m[2])
			pq.Options[opt] = val
			optionsStarted = true
			continue
		}

		// ── Question text (before options start) ─────────────────────────────
		if !optionsStarted {
			qLines = append(qLines, t)
		}
	}

	pq.Text = strings.Join(qLines, " ")
	pq.Explanation = strings.Join(explanationLines, " ")

	return pq
}
