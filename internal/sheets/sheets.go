package sheets

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"question_input_smartsystem/internal/formatter"
)

// ExportResult is returned after a successful export.
type ExportResult struct {
	SheetURL  string `json:"sheet_url"`
	SheetName string `json:"sheet_name"`
}

// Header columns for the Google Sheet (exact requested format).
var headers = []interface{}{
	"text", "difficulty", "module",
	"option1_text", "option1_is_correct",
	"option2_text", "option2_is_correct",
	"option3_text", "option3_is_correct",
	"option4_text", "option4_is_correct",
}

// Export creates a new Google Sheet, writes headers and all rows in one batch call.
// credentialsPath must point to a Google Service Account JSON key file.
func Export(credentialsPath string, rows []formatter.SheetRow, meta formatter.Metadata) (ExportResult, error) {
	ctx := context.Background()

	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credentialsPath), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		return ExportResult{}, fmt.Errorf("failed to create Sheets client: %w", err)
	}

	// --- 1. Build dynamic sheet name ---
	sheetName := buildSheetName(meta)

	// --- 2. Create the spreadsheet ---
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: sheetName,
		},
	}
	created, err := srv.Spreadsheets.Create(spreadsheet).Context(ctx).Do()
	if err != nil {
		return ExportResult{}, fmt.Errorf("failed to create spreadsheet: %w", err)
	}

	spreadsheetID := created.SpreadsheetId
	sheetTitle := created.Sheets[0].Properties.Title // default "Sheet1"

	// --- 3. Build all row data (header + data rows) ---
	values := make([]*sheets.ValueRange, 0, 1)

	var allValues [][]interface{}
	allValues = append(allValues, headers)

	for _, row := range rows {
		allValues = append(allValues, formatter.ToInterfaceSlice(row))
	}

	values = append(values, &sheets.ValueRange{
		Range:  fmt.Sprintf("%s!A1", sheetTitle),
		Values: allValues,
	})

	// --- 4. Batch update (single API call for all rows) ---
	batchReq := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
		Data:             values,
	}
	_, err = srv.Spreadsheets.Values.BatchUpdate(spreadsheetID, batchReq).Context(ctx).Do()
	if err != nil {
		return ExportResult{}, fmt.Errorf("failed to write data to spreadsheet: %w", err)
	}

	// --- 5. Bold the header row ---
	_, _ = srv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				RepeatCell: &sheets.RepeatCellRequest{
					Range: &sheets.GridRange{
						SheetId:          created.Sheets[0].Properties.SheetId,
						StartRowIndex:    0,
						EndRowIndex:      1,
						StartColumnIndex: 0,
						EndColumnIndex:   int64(len(headers)),
					},
					Cell: &sheets.CellData{
						UserEnteredFormat: &sheets.CellFormat{
							TextFormat: &sheets.TextFormat{Bold: true},
							BackgroundColor: &sheets.Color{
								Red: 0.20, Green: 0.60, Blue: 0.86,
							},
						},
					},
					Fields: "userEnteredFormat(textFormat,backgroundColor)",
				},
			},
		},
	}).Context(ctx).Do()

	return ExportResult{
		SheetURL:  fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", spreadsheetID),
		SheetName: sheetName,
	}, nil
}

// buildSheetName creates a name like: Class4_Mixed_2026_03_27_10AM
func buildSheetName(meta formatter.Metadata) string {
	now := time.Now()

	class := sanitize(meta.Class)
	module := sanitize(meta.Module)

	date := now.Format("2006_01_02")
	hour := now.Format("3PM") // e.g. "10AM"

	parts := []string{}
	if class != "" {
		parts = append(parts, class)
	}
	if module != "" {
		parts = append(parts, module)
	}
	parts = append(parts, date, hour)

	return strings.Join(parts, "_")
}

func sanitize(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "")
	// Remove non-alphanumeric characters
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}
