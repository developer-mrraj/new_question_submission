package handler

import (
	"encoding/json"
	"net/http"
)

// ParseRequest is the JSON body for POST /parse
type ParseRequest struct {
	Difficulty string          `json:"difficulty" example:"Easy"`
	Module     string          `json:"module" example:"Mixed"`
	Questions  []InputQuestion `json:"questions"`
}

type InputQuestion struct {
	Number  int               `json:"number"`
	Text    string            `json:"text"`
	Options map[string]string `json:"options"`
	Answer  string            `json:"answer"`
}

// ParsedFlatQuestion represents the flattened output format
type ParsedFlatQuestion struct {
	Text             string `json:"text"`
	Difficulty       string `json:"difficulty"`
	Module           string `json:"module"`
	Option1Text      string `json:"option1_text"`
	Option1IsCorrect bool   `json:"option1_is_correct"`
	Option2Text      string `json:"option2_text"`
	Option2IsCorrect bool   `json:"option2_is_correct"`
	Option3Text      string `json:"option3_text"`
	Option3IsCorrect bool   `json:"option3_is_correct"`
	Option4Text      string `json:"option4_text"`
	Option4IsCorrect bool   `json:"option4_is_correct"`
}

// ErrorResponse is a generic error envelope.
type ErrorResponse struct {
	Error string `json:"error" example:"Field 'text' is required"`
}

// ParseHandler handles POST /parse
//
//	@Summary		Parse structured MCQ JSON into flattened format
//	@Description	Accepts JSON structured questions, validates them, and flattens options A-D with correctness flags.
//	@Tags			Parser
//	@Accept			json
//	@Produce		json
//	@Param			body	body		ParseRequest	true	"Structured questions + metadata"
//	@Success		200		{array}		ParsedFlatQuestion
//	@Failure		400		{object}	ErrorResponse
//	@Router			/parse [post]
func ParseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	result := make([]ParsedFlatQuestion, 0)

	for _, q := range req.Questions {
		// Validation Rules
		if len(q.Options) != 4 {
			continue // Skip invalid (not exactly 4 options)
		}
		
		valA, okA := q.Options["A"]
		valB, okB := q.Options["B"]
		valC, okC := q.Options["C"]
		valD, okD := q.Options["D"]

		if !okA || !okB || !okC || !okD {
			continue // Skip invalid (missing A, B, C, or D)
		}

		if q.Answer != "A" && q.Answer != "B" && q.Answer != "C" && q.Answer != "D" {
			continue // Skip invalid (unknown answer)
		}

		flat := ParsedFlatQuestion{
			Text:             q.Text,
			Difficulty:       req.Difficulty,
			Module:           req.Module,
			Option1Text:      valA,
			Option1IsCorrect: q.Answer == "A",
			Option2Text:      valB,
			Option2IsCorrect: q.Answer == "B",
			Option3Text:      valC,
			Option3IsCorrect: q.Answer == "C",
			Option4Text:      valD,
			Option4IsCorrect: q.Answer == "D",
		}
		result = append(result, flat)
	}

	// Always return an array
	if result == nil {
		result = []ParsedFlatQuestion{}
	}

	SetLastParsedData(result)

	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}
