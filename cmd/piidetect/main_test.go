package main

import (
	"testing"

	"github.com/ullauri/piidetect/config"
)

func TestGetExpandedPaths(t *testing.T) {
	paths, err := getExpandedPaths([]string{"./../../testdata"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(paths) == 0 {
		t.Fatalf("expected paths, got empty list")
	}
}

func TestGetIssues(t *testing.T) {
	config.DefaultSetup()
	paths := []string{"./../../testdata/file.go"}
	issues, err := getIssues(paths)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(issues) == 0 {
		t.Fatalf("expected issues, got none")
	}
}
