package preprocessor_test

import (
	"strings"
	"testing"

	"question_input_smartsystem/internal/preprocessor"
)

func TestRemovesEmojis(t *testing.T) {
	// Line starting with 🧠 should be completely dropped
	input1 := "🧠 Section Header Test\nQ1. What is 2+2?\nA. 3\nB. 4\nC. 5\nD. 6\nAnswer: B"
	result1 := preprocessor.Process(input1)
	if strings.Contains(result1, "🧠 Section Header Test") || strings.Contains(result1, "🧠") {
		t.Errorf("Emoji header line was not removed entirely: %q", result1)
	}

	// Inline emojis should just be stripped
	input2 := "Q2. 🧪 What is water?\nA. H2O\nB. O2\nC. CO2\nD. N2\nAnswer: A"
	result2 := preprocessor.Process(input2)
	if strings.Contains(result2, "🧪") {
		t.Errorf("Inline emoji not stripped: %q", result2)
	}
	if !strings.Contains(result2, "What is water?") {
		t.Errorf("Inline text was altered improperly: %q", result2)
	}
}

func TestNormalizesLineEndings(t *testing.T) {
	input := "Q1. Test\r\nA. Yes\r\nB. No\r\nC. Maybe\r\nD. Never\r\nAnswer: A"
	result := preprocessor.Process(input)
	if strings.Contains(result, "\r") {
		t.Errorf("Carriage returns not removed: %q", result)
	}
}

func TestRemovesHeadings(t *testing.T) {
	input := "Class 4 Mixed Test\nQ1. What is 2+2?\nA. 3\nB. 4\nC. 5\nD. 6\nAnswer: B"
	result := preprocessor.Process(input)
	if strings.Contains(result, "Class 4 Mixed Test") {
		t.Errorf("Heading not removed: %q", result)
	}
	if !strings.Contains(result, "Q1.") {
		t.Errorf("Question line was removed: %q", result)
	}
}

func TestPreservesOptionLines(t *testing.T) {
	input := "Q1. Capital?\nA. Delhi\nB. Mumbai\nC. Chennai\nD. Pune\nAnswer: A"
	result := preprocessor.Process(input)
	for _, opt := range []string{"A. Delhi", "B. Mumbai", "C. Chennai", "D. Pune"} {
		if !strings.Contains(result, opt) {
			t.Errorf("Option line was altered, missing: %q in result: %q", opt, result)
		}
	}
}
