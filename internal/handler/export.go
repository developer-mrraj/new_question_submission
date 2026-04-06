package handler

import (
	"fmt"
	"net/http"

	"github.com/xuri/excelize/v2"
)

// ExportHandler handles GET /export
//
//	@Summary		Export parsed questions to Excel
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
