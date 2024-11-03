package detector

import (
	"context"
	"testing"

	"github.com/ullauri/piidetect/config"
)

func TestDetectPII(t *testing.T) {
	config.DefaultSetup()
	path := "./../testdata/file.go"
	issues, err := DetectPII(context.Background(), path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(issues) == 0 {
		t.Fatal("expected issues, but got none")
	}
}

func TestDetectPIIAST(t *testing.T) {
	config.DefaultSetup()
	path := "./../testdata/file.go"
	issues, err := detectPIIAST(context.Background(), path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(issues) == 0 {
		t.Fatal("expected issues, but got none")
	}
}

func TestIsGoFile(t *testing.T) {
	if !isGoFile("file.go") {
		t.Fatal("expected true for Go file")
	}
	if isGoFile("file.txt") {
		t.Fatal("expected false for non-Go file")
	}
}
