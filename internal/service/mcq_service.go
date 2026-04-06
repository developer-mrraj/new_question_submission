package service

import (
	"regexp"

	"question_input_smartsystem/internal/model"
	"question_input_smartsystem/internal/parser"
)

type MCQService struct{}

func NewMCQService() *MCQService {
	return &MCQService{}
}

// reHasAnswer detects "Answer:" or "✅ Answer:" (case-insensitive).
var reHasAnswer = regexp.MustCompile(`(?i)(?:✅\s*)?answer:`)

// reHasLevelOpt detects options written as "a)" or "A)" at the start of a line (case-insensitive).
var reHasLevelOpt = regexp.MustCompile(`(?im)^\s*[a-d]\)`)

func (s *MCQService) ConvertTextToJSON(text string) []model.MCQ {
	// Route to level parser when options use "a) / A)" style
	if reHasAnswer.MatchString(text) && reHasLevelOpt.MatchString(text) {
		return parser.ParseLevelMCQs(text)
	}

	// Default: standard A. / A) / A- format handled by ParseMCQs
	return parser.ParseMCQs(text)
}
