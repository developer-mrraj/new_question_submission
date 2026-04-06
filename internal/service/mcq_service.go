package service

import (
	"strings"

	"question_input_smartsystem/internal/model"
	"question_input_smartsystem/internal/parser"
)

type MCQService struct{}

func NewMCQService() *MCQService {
	return &MCQService{}
}

func (s *MCQService) ConvertTextToJSON(text string) []model.MCQ {
	// Route to level parser if input uses a) b) c) format
	if strings.Contains(text, "Answer:") && strings.Contains(text, "a)") {
		return parser.ParseLevelMCQs(text)
	}

	// Default: standard A. B. C. D. format
	return parser.ParseMCQs(text)
}
