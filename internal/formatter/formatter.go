package formatter

import (
	"question_input_smartsystem/internal/parser"
)

// SheetRow represents a single row in the Google Sheet.
type SheetRow struct {
	QuestionNo     int    `json:"question_no"`
	Text           string `json:"text"`
	Difficulty     string `json:"difficulty"`
	Module         string `json:"module"`
	Class          string `json:"class"`
	Option1Text    string `json:"option1_text"`
	Option1Correct bool   `json:"option1_is_correct"`
	Option2Text    string `json:"option2_text"`
	Option2Correct bool   `json:"option2_is_correct"`
	Option3Text    string `json:"option3_text"`
	Option3Correct bool   `json:"option3_is_correct"`
	Option4Text    string `json:"option4_text"`
	Option4Correct bool   `json:"option4_is_correct"`
}

// Metadata holds the class/module/difficulty supplied by the user.
//	@Description	Contextual metadata attached to the question set for export
type Metadata struct {
	Class      string `json:"class" example:"Class4"`
	Module     string `json:"module" example:"Mixed"`
	Difficulty string `json:"difficulty" example:"Easy"`
}

// Format converts a slice of valid ParsedQuestions into SheetRows.
// Options are mapped in fixed order A→1, B→2, C→3, D→4.
func Format(questions []parser.ParsedQuestion, meta Metadata) []SheetRow {
	optionOrder := []string{"A", "B", "C", "D"}
	rows := make([]SheetRow, 0, len(questions))

	for _, q := range questions {
		row := SheetRow{
			QuestionNo: q.Number,
			Text:       q.Text,
			Difficulty: meta.Difficulty,
			Module:     meta.Module,
			Class:      meta.Class,
		}

		texts := [4]string{}
		corrects := [4]bool{}
		for i, letter := range optionOrder {
			texts[i] = q.Options[letter]
			corrects[i] = (q.Answer == letter)
		}

		row.Option1Text = texts[0]
		row.Option1Correct = corrects[0]
		row.Option2Text = texts[1]
		row.Option2Correct = corrects[1]
		row.Option3Text = texts[2]
		row.Option3Correct = corrects[2]
		row.Option4Text = texts[3]
		row.Option4Correct = corrects[3]

		rows = append(rows, row)
	}
	return rows
}

// ToInterfaceSlice converts a SheetRow to []interface{} for the Sheets API.
// Matching requested order:
// text, difficulty, module, option1_text, option1_is_correct, option2_text, option2_is_correct, option3_text, option3_is_correct, option4_text, option4_is_correct
func ToInterfaceSlice(row SheetRow) []interface{} {
	correct := func(b bool) string {
		if b {
			return "TRUE"
		}
		return "FALSE"
	}
	return []interface{}{
		row.Text,
		row.Difficulty,
		row.Module,
		row.Option1Text,
		correct(row.Option1Correct),
		row.Option2Text,
		correct(row.Option2Correct),
		row.Option3Text,
		correct(row.Option3Correct),
		row.Option4Text,
		correct(row.Option4Correct),
	}
}
