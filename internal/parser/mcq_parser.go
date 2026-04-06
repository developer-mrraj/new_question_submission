package parser

import (
	"regexp"
	"strconv"
	"strings"

	"question_input_smartsystem/internal/model"
)

// reFirstOpt matches the first line that starts an option block (a. / a) / A. / A))
// anchored to start-of-line so "a)" inside question text is not a false positive.
var reFirstOptMCQ = regexp.MustCompile(`(?im)^\s*[a-d][\.\)]`)

// reOptionMCQ extracts a single option line anchored to start-of-line.
var reOptionMCQ = regexp.MustCompile(`(?im)^\s*([a-d])[\.\)]\s+(.+)`)

func ParseMCQs(text string) []model.MCQ {

	// Split by Q1, Q2... anchored to start of line
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
		// Use a line-anchored regex to find the first option line — avoids
		// false positives when the question body itself contains "a)" or "a.".
		loc := reFirstOptMCQ.FindStringIndex(block)
		optIndex := -1
		if loc != nil {
			optIndex = loc[0]
		}

		questionText := ""
		if optIndex != -1 {
			questionText = strings.TrimSpace(block[:optIndex])
		}

		// remove Q number prefix
		questionText = numRegex.ReplaceAllString(questionText, "")
		questionText = strings.ReplaceAll(questionText, "⸻", "")
		questionText = strings.TrimSpace(questionText)

		// 🧩 STEP 3: Extract Options
		// Remove answer line first so the answer letter is not captured as an option.
		cleanBlock := regexp.MustCompile(`(?i)(?:✅\s*)?answer:.*`).ReplaceAllString(block, "")

		// Anchored to start-of-line — prevents mid-sentence letter matches.
		optMatches := reOptionMCQ.FindAllStringSubmatch(cleanBlock, -1)

		options := make(map[string]string)
		for _, opt := range optMatches {
			key := strings.ToUpper(opt[1])
			value := strings.TrimSpace(opt[2])

			// take only the first line of the value
			value = strings.Split(value, "\n")[0]
			value = strings.TrimSpace(value)

			options[key] = value
		}

		// 🎯 STEP 4: Extract Answer (UNIVERSAL)
		// Matches: ✅ Answer: B, Answer: B, ✅ Ans: B, Ans: B
		// Also handles "Ans: B) 39" where letter is followed by ) and text.
		ansRegex := regexp.MustCompile(`(?i)(?:✅\s*)?(?:answer|ans):\s*([a-dA-D])`)
		ansMatches := ansRegex.FindAllStringSubmatch(block, -1)

		answer := ""
		if len(ansMatches) > 0 {
			// take LAST match to avoid spurious early matches
			answer = strings.ToUpper(ansMatches[len(ansMatches)-1][1])
		}

		// 🧠 FALLBACK: match answer text against option values
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

		// ✅ Validation: require all 4 options (A-D) explicitly
		_, hasA := options["A"]
		_, hasB := options["B"]
		_, hasC := options["C"]
		_, hasD := options["D"]
		if hasA && hasB && hasC && hasD {
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
