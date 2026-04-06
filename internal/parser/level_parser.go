package parser

import (
	"regexp"
	"strconv"
	"strings"

	"question_input_smartsystem/internal/model"
)

// ParseLevelMCQs parses MCQ text that uses a) b) c) d) option style
// and "Answer: X" or "✅ Answer: X)" answer lines.
// Handles multi-line questions (e.g. fill-in-the-blank across multiple lines).
func ParseLevelMCQs(text string) []model.MCQ {
	// ✅ FIXED: improved block regex (no dependency on Answer)
	blockRegex := regexp.MustCompile(`(?s)(\d+[\.\)]\s*.*?)(?=\n\d+[\.\)]|\z)`)
	blocks := blockRegex.FindAllString(text, -1)

	var results []model.MCQ

	// option pattern: a) some text  OR  a. some text
	optRegex := regexp.MustCompile(`(?i)^([a-d])[\)\.]\s*(.+)$`)

	// ✅ FIXED: universal answer regex
	ansRegex := regexp.MustCompile(`(?i)(?:✅\s*)?answer:\s*([a-d])`)

	// question number prefix: "1. " or "1) "
	numRegex := regexp.MustCompile(`^(\d+)[\.\)]\s*`)

	// separator cleanup
	separatorRegex := regexp.MustCompile(`[⸻—–-]{2,}`)

	for _, block := range blocks {
		block = strings.TrimSpace(block)

		// ── Fix 1: extract question number from the very first token ──
		numMatch := numRegex.FindStringSubmatch(block)
		qNum := 0
		if len(numMatch) > 1 {
			qNum, _ = strconv.Atoi(numMatch[1])
		}

		// ── Fix 2: extract answer from FULL BLOCK (not line-by-line) ──
		ansMatches := ansRegex.FindAllStringSubmatch(block, -1)
		answer := ""
		if len(ansMatches) > 0 {
			// take LAST match (important)
			answer = strings.ToUpper(ansMatches[len(ansMatches)-1][1])
		}

		// ── Fix 3: remove answer before parsing options ──
		cleanBlock := regexp.MustCompile(`(?i)(?:✅\s*)?answer:.*`).ReplaceAllString(block, "")

		// ── Fix 4: find first option position (multi-line safe) ──
		lowerBlock := strings.ToLower(cleanBlock)
		optIndex := strings.Index(lowerBlock, "\na)")
		if optIndex == -1 {
			optIndex = strings.Index(lowerBlock, "a)")
		}
		if optIndex == -1 {
			optIndex = strings.Index(lowerBlock, "\na.")
		}
		if optIndex == -1 {
			optIndex = strings.Index(lowerBlock, "a.")
		}

		if optIndex == -1 {
			continue
		}

		// Everything before the first option = question text
		questionText := cleanBlock[:optIndex]

		// Strip leading "1. " or "1) "
		questionText = numRegex.ReplaceAllString(questionText, "")

		// Strip decorative separators
		questionText = separatorRegex.ReplaceAllString(questionText, "")

		questionText = strings.TrimSpace(questionText)

		// ── Parse options ──
		rest := cleanBlock[optIndex:]
		lines := strings.Split(rest, "\n")

		options := make(map[string]string)

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Option line
			if optMatch := optRegex.FindStringSubmatch(line); len(optMatch) >= 3 {
				key := strings.ToUpper(optMatch[1])
				val := strings.TrimSpace(optMatch[2])
				options[key] = val
			}
		}

		// 🧪 DEBUG (optional but useful)
		// if answer == "" {
		// 	fmt.Println("⚠️ Missing answer:", block)
		// }

		// Validation
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
