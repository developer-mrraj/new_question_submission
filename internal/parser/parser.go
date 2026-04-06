package parser

import (
	"regexp"
	"strings"
)

type ParsedQuestion struct {
	Number  int               `json:"number"`
	Text    string            `json:"text"`
	Options map[string]string `json:"options"`
	Answer  string            `json:"answer"`
}

// reQHeader strips the prefix "Q1.", "1)", "🔥 Q19", etc. from the first line.
var reQHeader = regexp.MustCompile(`(?i)^\s*(?:🔥\s*)?(?:q\s*\d+(?:\s*[\.\)])?|\d+\s*[\.\)])\s*`)

// reOption supports: A) A. A- and even "A - " (space-dash-space) formats.
var reOption = regexp.MustCompile(`(?i)^\s*([a-d])(?:[\)\.]|\s*-)\s*(.+)`)

// reAnswer requires the keyword 'answer' or 'ans', followed by ANY chars up to a colon,
// then the option character. This handles: "Answer: c)" as well as "Answer (including pets): c)".
var reAnswer = regexp.MustCompile(`(?i)(?:answer|ans).*?:\s*([a-d])`)

func Parse(num int, raw string) ParsedQuestion {
	lines := strings.Split(raw, "\n")

	pq := ParsedQuestion{
		Number:  num,
		Options: map[string]string{},
	}

	var qLines []string
	optionsStarted := false

	for i, line := range lines {
		t := strings.TrimSpace(line)

		if t == "" {
			continue
		}

		// remove Q1
		if i == 0 {
			t = reQHeader.ReplaceAllString(t, "")
			qLines = append(qLines, t)
			continue
		}

		// answer
		if m := reAnswer.FindStringSubmatch(t); m != nil {
			pq.Answer = strings.ToUpper(m[1])
			continue
		}

		// options (a), b), A., etc)
		if m := reOption.FindStringSubmatch(t); m != nil {
			opt := strings.ToUpper(m[1])
			val := strings.TrimSpace(m[2])

			pq.Options[opt] = val
			optionsStarted = true
			continue
		}

		// question text
		if !optionsStarted {
			qLines = append(qLines, t)
		}
	}

	pq.Text = strings.Join(qLines, " ")

	return pq
}
