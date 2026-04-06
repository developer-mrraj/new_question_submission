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

var reQHeader = regexp.MustCompile(`(?i)^\s*(?:q\s*)?\d+\s*[\.\)]\s*`)
var reOption = regexp.MustCompile(`(?i)^\s*([a-d])[\)\.\-]\s+(.+)`)
var reAnswer = regexp.MustCompile(`(?i)(?:answer|ans|correct).*?([a-d])`)

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

			// ❌ skip garbage options
			if len(val) < 2 {
				continue
			}

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
