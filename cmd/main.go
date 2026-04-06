package main

import (
	"log"
	"net/http"

	"question_input_smartsystem/internal/handler"
	"question_input_smartsystem/internal/service"

	_ "question_input_smartsystem/docs" // swagger docs
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title MCQ Parser API
// @version 1.0
// @description API to parse raw text into structured MCQ JSON
// @host localhost:8080
// @BasePath /
func main() {

	service := service.NewMCQService()
	handler := handler.NewMCQHandler(service)

	http.HandleFunc("/convert", handler.ConvertHandler)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
