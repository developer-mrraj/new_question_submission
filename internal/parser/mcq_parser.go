package parser

import (
	"regexp"
	"strconv"
	"strings"

	"question_input_smartsystem/internal/model"
)

func ParseMCQs(text string) []model.MCQ {

	// Split by Q1, Q2... (more reliable than old regex)
	qRegex := regexp.MustCompile(`(?m)^Q\d+\.`)
	indices := qRegex.FindAllStringIndex(text, -1)

	var results []model.MCQ

	for i := 0; i < len(indices); i++ {

		start := indices[i][0]
		var end int

		if i+1 < len(indices) {
			end = indices[i+1][0]
		} else {
			end = len(text)
		}

		block := strings.TrimSpace(text[start:end])

		// 🔥 STEP 1: Clean unwanted parts
		block = strings.Split(block, "Explanation:")[0]
		block = strings.Split(block, "👉")[0]

		// 🔢 Extract Question Number
		numRegex := regexp.MustCompile(`Q(\d+)\.`)
		numMatch := numRegex.FindStringSubmatch(block)

		if len(numMatch) < 2 {
			continue
		}

		qNum, _ := strconv.Atoi(numMatch[1])

		// 🧠 STEP 2: Extract Question Text (MULTILINE SAFE)
		optIndex := strings.Index(strings.ToLower(block), "\na)")
		if optIndex == -1 {
			optIndex = strings.Index(strings.ToLower(block), "a)")
		}
		if optIndex == -1 {
			optIndex = strings.Index(strings.ToLower(block), "\na.")
		}
		if optIndex == -1 {
			optIndex = strings.Index(strings.ToLower(block), "a.")
		}

		questionText := ""
		if optIndex != -1 {
			questionText = strings.TrimSpace(block[:optIndex])
		}

		// remove Q number
		questionText = numRegex.ReplaceAllString(questionText, "")
		questionText = strings.ReplaceAll(questionText, "⸻", "")
		questionText = strings.TrimSpace(questionText)

		// 🧩 STEP 3: Extract Options (supports A. and a))
		options := make(map[string]string)

		optRegex := regexp.MustCompile(`(?i)([a-d])[\.\)]\s*(.+)`)

		// ✅ FIX: remove answer before extracting options
		cleanBlock := regexp.MustCompile(`(?i)(?:✅\s*)?answer:.*`).ReplaceAllString(block, "")

		optMatches := optRegex.FindAllStringSubmatch(cleanBlock, -1)

		for _, opt := range optMatches {
			key := strings.ToUpper(opt[1])
			value := strings.TrimSpace(opt[2])

			// take only first line
			value = strings.Split(value, "\n")[0]

			options[key] = value
		}

		// 🎯 STEP 4: Extract Answer (UNIVERSAL)

		// ❌ old:
		// ansRegex := regexp.MustCompile(`(?i)answer:\s*([a-dA-D])`)
		// ansMatch := ansRegex.FindStringSubmatch(block)

		// answer := ""
		// if len(ansMatch) > 1 {
		// 	answer = strings.ToUpper(ansMatch[1])
		// }

		// ✅ FIXED:
		ansRegex := regexp.MustCompile(`(?i)(?:✅\s*)?answer:\s*([a-dA-D])`)
		ansMatches := ansRegex.FindAllStringSubmatch(block, -1)

		answer := ""
		if len(ansMatches) > 0 {
			// take LAST match (important)
			answer = strings.ToUpper(ansMatches[len(ansMatches)-1][1])
		}

		// 🧪 DEBUG (important)
		// 🧠 FALLBACK: try to detect answer from option text
		if answer == "" {

			ansTextRegex := regexp.MustCompile(`(?i)(?:answer|ans).*?:?\s*([^\n]+)`)
			textMatch := ansTextRegex.FindStringSubmatch(block)

			if len(textMatch) > 1 {
				ansText := strings.TrimSpace(textMatch[1])

				for key, val := range options {
					if strings.Contains(strings.ToLower(ansText), strings.ToLower(val)) {
						answer = key
						break
					}
				}
			}
		}
		// ✅ Validation (relaxed)
		if len(options) >= 4 {
			results = append(results, model.MCQ{
				QuestionNumber: qNum,
				Title:          "",
				Question:       questionText,
				Options:        options,
				CorrectAnswer:  answer,
			})
		}
	}

	return results
}
