package config

import (
	"os"
	"testing"
	"time"

	"github.com/ullauri/piidetect"
)

func TestSetup(t *testing.T) {
	Setup(
		WithMethod(piidetect.AST),
		WithTotalWorkers(5),
		WithTimeout(15*time.Second),
		WithPIIPatterns([]string{"ssn", "email"}),
	)

	if Method() != piidetect.AST {
		t.Fatalf("expected AST method, got %v", Method())
	}
	if TotalWorkers() != 5 {
		t.Fatalf("expected 5 workers, got %v", TotalWorkers())
	}
	if Timeout() != 15*time.Second {
		t.Fatalf("expected 15s timeout, got %v", Timeout())
	}
	if len(PIIPatterns()) != 2 {
		t.Fatalf("expected 2 PII patterns, got %d", len(PIIPatterns()))
	}
}

func TestReadPIIPatterns(t *testing.T) {
	filePath := "./../testdata/patterns.txt"
	err := os.WriteFile(filePath, []byte("email\nphone\n"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer os.Remove(filePath)

	patterns, err := ReadPIIPatterns(filePath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := 2
	if len(patterns) != expected {
		t.Fatalf("expected %d patterns, got %d", expected, len(patterns))
	}
}
