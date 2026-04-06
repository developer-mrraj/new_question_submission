// diagnose/main.go — run with: go run ./cmd/diagnose/main.go <inputfile.txt>
// Prints every block from the splitter with VALID/INVALID verdict and the raw block.
package main

import (
	"fmt"
	"io"
	"os"

	"question_input_smartsystem/internal/parser"
	"question_input_smartsystem/internal/preprocessor"
	"question_input_smartsystem/internal/splitter"
)

func main() {
	path := "test_payload_user.txt"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot open file:", err)
		os.Exit(1)
	}
	defer f.Close()

	b, _ := io.ReadAll(f)
	reqText := string(b)

	cleaned := preprocessor.Process(reqText)
	blocks := splitter.Split(cleaned)

	fmt.Printf("=== TOTAL BLOCKS FROM SPLITTER: %d ===\n\n", len(blocks))

	valid := 0
	invalid := 0
	for _, block := range blocks {
		q := parser.Parse(block.Number, block.Raw)
		missingOpts := []string{}
		for _, opt := range []string{"A", "B", "C", "D"} {
			if _, ok := q.Options[opt]; !ok {
				missingOpts = append(missingOpts, opt)
			}
		}

		if q.Text != "" && len(q.Options) >= 4 && q.Answer != "" && len(missingOpts) == 0 {
			valid++
			fmt.Printf("[VALID]   Block#%d  Q%-3d  options=%d  ans=%s  text=%q\n",
				block.Number, q.Number, len(q.Options), q.Answer, truncate(q.Text, 60))
		} else {
			invalid++
			reasons := []string{}
			if q.Text == "" {
				reasons = append(reasons, "NO TEXT")
			}
			if q.Answer == "" {
				reasons = append(reasons, "NO ANSWER")
			}
			if len(missingOpts) > 0 {
				reasons = append(reasons, fmt.Sprintf("MISSING OPTIONS: %v", missingOpts))
			}
			fmt.Printf("[INVALID] Block#%d  Q%-3d  REASONS: %v\n", block.Number, q.Number, reasons)
			fmt.Printf("  RAW BLOCK >>>\n%s\n  <<<\n\n", indent(block.Raw))
		}
	}

	fmt.Printf("\n=== SUMMARY: %d valid / %d invalid / %d total ===\n", valid, invalid, len(blocks))
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "…"
}

func indent(s string) string {
	lines := ""
	for _, l := range splitLines(s) {
		lines += "    " + l + "\n"
	}
	return lines
}

func splitLines(s string) []string {
	result := []string{}
	cur := ""
	for _, c := range s {
		if c == '\n' {
			result = append(result, cur)
			cur = ""
		} else {
			cur += string(c)
		}
	}
	if cur != "" {
		result = append(result, cur)
	}
	return result
}
