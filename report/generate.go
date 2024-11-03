package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ullauri/piidetect"
	"github.com/ullauri/piidetect/config"
)

func Generate(issues []piidetect.Issue) {
	extension := filepath.Ext(config.OutputFilePath())
	if extension == "" {
		generateText(issues)
	} else if extension == ".json" {
		err := generateJSON(issues)
		if err != nil {
			fmt.Printf("error generating JSON report: %v\n", err)
		}
	} else {
		fmt.Printf("unsupported file extension: %s\n", extension)
		generateText(issues)
	}
}

func generateJSON(issues []piidetect.Issue) error {
	filename := config.OutputFilePath()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(issues)
}

func generateText(issues []piidetect.Issue) {
	fmt.Print("potential PII Issues:\n\n")
	for _, issue := range issues {
		fmt.Printf("(%s) (%s) %s:%d: %s\n", issue.Match, issue.Type, issue.File, issue.Line, issue.Message)
	}
	fmt.Printf("\ntotal Issues: %d\n", len(issues))
}
