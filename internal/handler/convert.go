package handler

import (
	"io"
	"net/http"

	"question_input_smartsystem/internal/parser"
	"question_input_smartsystem/internal/preprocessor"
	"question_input_smartsystem/internal/splitter"
)

// ConvertHandler godoc
//
//	@Summary		Convert raw text into structured MCQ JSON
//	@Description	Paste raw MCQ text. Cleans, splits, parses and returns ALL questions (no filtering).
//	@Tags			Converter
//	@Accept			plain
//	@Produce		json
//	@Param			text	body		string	true	"Raw question text"
//	@Success		200		{array}	parser.ParsedQuestion
//	@Failure		400		{object}	ErrorResponse
//	@Router			/convert [post]
func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to read body")
		return
	}
	defer r.Body.Close()

	input := string(body)
	if input == "" {
		writeError(w, http.StatusBadRequest, "Empty body")
		return
	}

	// ✅ Step 1: Clean
	cleaned := preprocessor.Process(input)

	// ✅ Step 2: Split
	blocks := splitter.Split(cleaned)

	parsed := make([]parser.ParsedQuestion, 0)

	// ✅ Step 3: Parse (NO FILTERING)
	for _, b := range blocks {
		q := parser.Parse(b.Number, b.Raw)

		// 🔥 NO CONDITIONS — KEEP EVERYTHING
		parsed = append(parsed, q)
	}

	// ✅ Step 4: Return ALL
	writeJSON(w, http.StatusOK, parsed)
}
