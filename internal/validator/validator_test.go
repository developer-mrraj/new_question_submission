package validator_test

import (
	"testing"

	"question_input_smartsystem/internal/parser"
	"question_input_smartsystem/internal/validator"
)

func makeQ(num int, text, answer string, opts map[string]string) parser.ParsedQuestion {
	return parser.ParsedQuestion{Number: num, Text: text, Answer: answer, Options: opts}
}

var fullOpts = map[string]string{"A": "a", "B": "b", "C": "c", "D": "d"}

func TestValidQuestionPasses(t *testing.T) {
	q := makeQ(1, "What is 2+2?", "B", fullOpts)
	valid, errs := validator.Validate([]parser.ParsedQuestion{q})
	if len(errs) != 0 {
		t.Errorf("Expected no errors, got: %v", errs)
	}
	if len(valid) != 1 {
		t.Errorf("Expected 1 valid question, got %d", len(valid))
	}
}

func TestMissingOptionDetected(t *testing.T) {
	opts := map[string]string{"A": "a", "B": "b", "D": "d"} // missing C
	q := makeQ(2, "Some question?", "A", opts)
	_, errs := validator.Validate([]parser.ParsedQuestion{q})
	if len(errs) == 0 {
		t.Fatal("Expected validation error for missing option C")
	}
	if errs[0].QuestionNo != 2 {
		t.Errorf("Expected question_no=2, got %d", errs[0].QuestionNo)
	}
}

func TestMissingAnswerDetected(t *testing.T) {
	q := makeQ(3, "Some question?", "", fullOpts)
	_, errs := validator.Validate([]parser.ParsedQuestion{q})
	if len(errs) == 0 {
		t.Fatal("Expected validation error for missing answer")
	}
}

func TestEmptyTextDetected(t *testing.T) {
	q := makeQ(4, "", "A", fullOpts)
	_, errs := validator.Validate([]parser.ParsedQuestion{q})
	if len(errs) == 0 {
		t.Fatal("Expected validation error for empty text")
	}
}

func TestNeverStopsOnError(t *testing.T) {
	// Mix of valid and invalid
	q1 := makeQ(1, "Valid?", "A", fullOpts)
	q2 := makeQ(2, "Invalid - no answer", "", fullOpts)
	q3 := makeQ(3, "Also valid?", "C", fullOpts)

	valid, errs := validator.Validate([]parser.ParsedQuestion{q1, q2, q3})
	if len(valid) != 2 {
		t.Errorf("Expected 2 valid questions, got %d", len(valid))
	}
	if len(errs) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errs))
	}
}
