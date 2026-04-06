package splitter_test

import (
	"testing"

	"question_input_smartsystem/internal/splitter"
)

func TestSplitsTwoQuestions(t *testing.T) {
	input := "Q1. What is 2+2?\nA. 3\nB. 4\nC. 5\nD. 6\nAnswer: B\nQ2. Capital of India?\nA. Mumbai\nB. Delhi\nC. Chennai\nD. Pune\nAnswer: B"
	blocks := splitter.Split(input)
	if len(blocks) != 2 {
		t.Fatalf("Expected 2 blocks, got %d", len(blocks))
	}
	if blocks[0].Number != 1 {
		t.Errorf("Block 0 should have number=1, got %d", blocks[0].Number)
	}
	if blocks[1].Number != 2 {
		t.Errorf("Block 1 should have number=2, got %d", blocks[1].Number)
	}
}

func TestSplitsWithVariousFormats(t *testing.T) {
	// strict Q1. format
	input := "Q1. Question one\nA. a\nB. b\nC. c\nD. d\nAnswer: A\nQ2. Question two\nA. a\nB. b\nC. c\nD. d\nAnswer: B"
	blocks := splitter.Split(input)
	if len(blocks) != 2 {
		t.Fatalf("Expected 2 blocks, got %d", len(blocks))
	}
}

func TestIgnoresTextBeforeFirstQuestion(t *testing.T) {
	input := "Some heading text\nMore intro text\nQ1. Real question\nA. a\nB. b\nC. c\nD. d\nAnswer: A"
	blocks := splitter.Split(input)
	if len(blocks) != 1 {
		t.Fatalf("Expected 1 block, got %d", len(blocks))
	}
}
