package model

type MCQ struct {
	QuestionNumber int               `json:"question_number"`
	Title          string            `json:"title"`
	Question       string            `json:"question"`
	Options        map[string]string `json:"options"`
	CorrectAnswer  string            `json:"correct_answer"`
}
