package parser

import (
	"regexp"
	"strconv"
	"strings"

	"question_input_smartsystem/internal/model"
)

// ParseLevelMCQs parses MCQ text that uses a) b) c) d) / A) B) C) D) option style
// and "Answer: X" or "✅ Answer: X" answer lines.
// Handles multi-line questions (e.g. fill-in-the-blank across multiple lines).
func ParseLevelMCQs(text string) []model.MCQ {

	// ✅ FIX: use a line-anchored regex so numbers inside sentences/answers
	// don't break block boundaries (e.g. "Answer: 2. some text" no longer splits).
	numLineRegex := regexp.MustCompile(`(?m)^\d+[\.\)]\s`)
	indices := numLineRegex.FindAllStringIndex(text, -1)

	var blocks []string
	for i, loc := range indices {
		start := loc[0]
		var end int
		if i+1 < len(indices) {
			end = indices[i+1][0]
		} else {
			end = len(text)
		}
		blocks = append(blocks, strings.TrimSpace(text[start:end]))
	}

	var results []model.MCQ

	// option pattern: a) some text  OR  a. some text — anchored to line start
	optRegex := regexp.MustCompile(`(?im)^\s*([a-d])[\)\.][ \t]+(.+)$`)

	// ✅ universal answer regex (case-insensitive, emoji-safe, handles "Ans:" shorthand)
	ansRegex := regexp.MustCompile(`(?i)(?:✅\s*)?(?:answer|ans):\s*([a-d])`)

	// question number prefix: "1. " or "1) "
	numRegex := regexp.MustCompile(`^(\d+)[\.\)]\s*`)

	// separator cleanup
	separatorRegex := regexp.MustCompile(`[⸻—–-]{2,}`)

	// ✅ FIX: line-anchored first-option detector — prevents false positives
	// when the question text itself contains "a)" or "a.".
	firstOptRegex := regexp.MustCompile(`(?im)^\s*[a-d][\)\.]`)

	for _, block := range blocks {
		block = strings.TrimSpace(block)

		// ── Extract question number ──
		numMatch := numRegex.FindStringSubmatch(block)
		qNum := 0
		if len(numMatch) > 1 {
			qNum, _ = strconv.Atoi(numMatch[1])
		}
		// Skip silently if no question number found
		if qNum == 0 {
			continue
		}

		// ── Extract answer from FULL BLOCK ──
		ansMatches := ansRegex.FindAllStringSubmatch(block, -1)
		answer := ""
		if len(ansMatches) > 0 {
			// take LAST match (handles duplicate patterns)
			answer = strings.ToUpper(ansMatches[len(ansMatches)-1][1])
		}

		// ── Remove answer line before parsing options ──
		cleanBlock := regexp.MustCompile(`(?i)(?:✅\s*)?answer:.*`).ReplaceAllString(block, "")

		// ── ✅ FIX: find first option position using line-anchored regex ──
		loc := firstOptRegex.FindStringIndex(cleanBlock)
		if loc == nil {
			continue
		}
		optIndex := loc[0]

		// Everything before the first option = question text
		questionText := cleanBlock[:optIndex]

		// Strip leading "1. " or "1) "
		questionText = numRegex.ReplaceAllString(questionText, "")

		// Strip decorative separators
		questionText = separatorRegex.ReplaceAllString(questionText, "")

		questionText = strings.TrimSpace(questionText)

		// ── Parse options ──
		rest := cleanBlock[optIndex:]
		optMatches := optRegex.FindAllStringSubmatch(rest, -1)

		options := make(map[string]string)
		for _, optMatch := range optMatches {
			key := strings.ToUpper(optMatch[1])
			val := strings.TrimSpace(optMatch[2])
			options[key] = val
		}

		// Validation: all 4 options, non-empty answer, non-empty question
		if len(options) != 4 || answer == "" || questionText == "" {
			continue
		}

		results = append(results, model.MCQ{
			QuestionNumber: qNum,
			Title:          "",
			Question:       questionText,
			Options:        options,
			CorrectAnswer:  answer,
		})
	}

	return results
}
