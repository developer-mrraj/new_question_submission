package main

import (
	"encoding/json"
	"fmt"
	"os"

	"question_input_smartsystem/internal/parser"
	"question_input_smartsystem/internal/preprocessor"
	"question_input_smartsystem/internal/splitter"
)

func main() {
	b, err := os.ReadFile("level17.txt")
	if err != nil {
		panic(err)
	}

	raw := string(b)
	cleaned := preprocessor.Process(raw)
	blocks := splitter.Split(cleaned)

	var questions []parser.ParsedQuestion
	for _, b := range blocks {
		q := parser.Parse(b.Number, b.Raw)
		questions = append(questions, q)
	}

	out, _ := json.MarshalIndent(questions, "", "  ")
	fmt.Println(string(out))
}
