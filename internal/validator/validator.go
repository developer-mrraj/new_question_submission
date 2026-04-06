package validator

import (
	"fmt"
	"strings"

	"question_input_smartsystem/internal/parser"
)

type ValidationError struct {
	QuestionNo int    `json:"question_no"`
	Error      string `json:"error"`
}

var validAnswers = map[string]bool{"A": true, "B": true, "C": true, "D": true}
var requiredOptions = []string{"A", "B", "C", "D"}

func Validate(questions []parser.ParsedQuestion) (valid []parser.ParsedQuestion, errors []ValidationError) {
	for _, q := range questions {
		var errs []string

		if strings.TrimSpace(q.Text) == "" {
			errs = append(errs, "Question text is empty")
		}

		for _, opt := range requiredOptions {
			if _, ok := q.Options[opt]; !ok {
				errs = append(errs, fmt.Sprintf("Missing option %s", opt))
			}
		}

		if q.Answer == "" {
			errs = append(errs, "Missing answer")
		} else if !validAnswers[q.Answer] {
			errs = append(errs, fmt.Sprintf("Invalid answer: %s", q.Answer))
		}

		if len(errs) > 0 {
			errors = append(errors, ValidationError{
				QuestionNo: q.Number,
				Error:      strings.Join(errs, "; "),
			})
		} else {
			valid = append(valid, q)
		}
	}
	return
}
