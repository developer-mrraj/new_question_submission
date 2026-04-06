package parser_test

import (
	"testing"

	"question_input_smartsystem/internal/parser"
)

func TestParseBasicDotFormat(t *testing.T) {
	raw := "Q1. What is 2+2?\nA. 3\nB. 4\nC. 5\nD. 6\nAnswer: B"
	q := parser.Parse(1, raw)

	if q.Text != "What is 2+2?" {
		t.Errorf("Text mismatch: %q", q.Text)
	}
	if q.Options["A"] != "3" {
		t.Errorf("Option A mismatch: %q", q.Options["A"])
	}
	if q.Options["B"] != "4" {
		t.Errorf("Option B mismatch: %q", q.Options["B"])
	}
	if q.Answer != "B" {
		t.Errorf("Answer mismatch: %q", q.Answer)
	}
}

func TestParseParenthesisFormat(t *testing.T) {
	raw := "Q2. Capital of India?\nA) Mumbai\nB) Delhi\nC) Chennai\nD) Pune\nAns: B"
	q := parser.Parse(2, raw)

	if q.Answer != "B" {
		t.Errorf("Answer mismatch: %q", q.Answer)
	}
	if q.Options["A"] != "Mumbai" {
		t.Errorf("Option A mismatch: %q", q.Options["A"])
	}
}

func TestParseDashFormat(t *testing.T) {
	raw := "Q3. Largest planet?\nA - Jupiter\nB - Saturn\nC - Neptune\nD - Uranus\nAnswer: A"
	q := parser.Parse(3, raw)

	if q.Options["A"] != "Jupiter" {
		t.Errorf("Option A (dash format) mismatch: %q", q.Options["A"])
	}
	if q.Answer != "A" {
		t.Errorf("Answer mismatch: %q", q.Answer)
	}
}

func TestParseAnswerWithDotAndText(t *testing.T) {
	raw := "Q4. Value of pi?\nA. 3\nB. 3.14\nC. 3.16\nD. 3.18\nAnswer: B. 3.14"
	q := parser.Parse(4, raw)
	// Answer line regex only picks up the letter
	if q.Answer != "B" {
		t.Errorf("Answer mismatch: %q", q.Answer)
	}
}

func TestParseMultilineQuestion(t *testing.T) {
	raw := "Q5. What is the\ncorrect answer?\nA. Yes\nB. No\nC. Maybe\nD. None\nAnswer: C"
	q := parser.Parse(5, raw)
	if q.Text != "What is the correct answer?" {
		t.Errorf("Multi-line question text mismatch: %q", q.Text)
	}
}
