package handler

import (
	"fmt"
	"net/http"

	"github.com/xuri/excelize/v2"
)

// ExportHandler handles GET /export
//
//	@Summary		Export parsed MCQ questions to Excel
//	@Description	Takes the recently parsed data from memory and returns an Excel file (.xlsx) download.
//	@Tags			Export
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Success		200		{string}	string	"questions.xlsx"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/export [get]
func ExportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := GetLastParsedData()
	if len(data) == 0 {
		writeError(w, http.StatusBadRequest, "No data available to export")
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Sheet1"
	headers := []string{
		"text", "difficulty", "module",
		"option1_text", "option1_is_correct",
		"option2_text", "option2_is_correct",
		"option3_text", "option3_is_correct",
		"option4_text", "option4_is_correct",
	}

	// Write headers
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Write data
	for i, q := range data {
		row := i + 2

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), q.Text)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), q.Difficulty)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), q.Module)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), q.Option1Text)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), q.Option1IsCorrect)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), q.Option2Text)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), q.Option2IsCorrect)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), q.Option3Text)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), q.Option3IsCorrect)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), q.Option4Text)
		f.SetCellValue(sheet, fmt.Sprintf("K%d", row), q.Option4IsCorrect)
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=\"questions.xlsx\"")

	if err := f.Write(w); err != nil {
		return
	}
}

// ExportDQHandler handles GET /export/dq
//
//	@Summary		Export parsed Decision & Quiz questions to Excel
//	@Description	Takes the recently parsed Decision & Quiz data from memory and returns an Excel file (.xlsx) with title and explanation columns.
//	@Tags			Decision & Quiz
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Success		200		{string}	string	"dq_questions.xlsx"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/export/dq [get]
func ExportDQHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := GetLastParsedDQData()
	if len(data) == 0 {
		writeError(w, http.StatusBadRequest, "No Decision & Quiz data available to export")
		return
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := "Sheet1"
	// 13 columns: title, text, explanation + standard 10
	headers := []string{
		"title", "text", "explanation", "difficulty", "module",
		"option1_text", "option1_is_correct",
		"option2_text", "option2_is_correct",
		"option3_text", "option3_is_correct",
		"option4_text", "option4_is_correct",
	}

	// Write headers
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Write data rows
	// Columns: A=title, B=text, C=explanation, D=difficulty, E=module,
	//          F=opt1_text, G=opt1_correct, H=opt2_text, I=opt2_correct,
	//          J=opt3_text, K=opt3_correct, L=opt4_text, M=opt4_correct
	for i, q := range data {
		row := i + 2

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), q.Title)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), q.Text)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), q.Explanation)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), q.Difficulty)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), q.Module)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), q.Option1Text)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), q.Option1IsCorrect)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), q.Option2Text)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), q.Option2IsCorrect)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), q.Option3Text)
		f.SetCellValue(sheet, fmt.Sprintf("K%d", row), q.Option3IsCorrect)
		f.SetCellValue(sheet, fmt.Sprintf("L%d", row), q.Option4Text)
		f.SetCellValue(sheet, fmt.Sprintf("M%d", row), q.Option4IsCorrect)
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=\"dq_questions.xlsx\"")

	if err := f.Write(w); err != nil {
		return
	}
}
