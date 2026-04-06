package handler

import (
	"encoding/json"
	"io"
	"net/http"

	_ "question_input_smartsystem/internal/model"
	"question_input_smartsystem/internal/service"
)

type MCQHandler struct {
	service *service.MCQService
}

func NewMCQHandler(s *service.MCQService) *MCQHandler {
	return &MCQHandler{service: s}
}

// ConvertHandler godoc
// @Summary Convert raw text to JSON MCQs
// @Description Parses raw text containing MCQs and converts it into a structured JSON array
// @Tags mcq
// @Accept text/plain
// @Produce json
// @Param text body string true "Raw text containing questions"
// @Success 200 {array} model.MCQ
// @Failure 400 {string} string "Invalid request"
// @Router /convert [post]
func (h *MCQHandler) ConvertHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	text := string(body)

	result := h.service.ConvertTextToJSON(text)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
