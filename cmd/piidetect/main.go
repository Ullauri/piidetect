package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ullauri/piidetect"
	"github.com/ullauri/piidetect/config"
	"github.com/ullauri/piidetect/detector"
	"github.com/ullauri/piidetect/report"
)

func main() {
	method := flag.String("method", "ast", "Detection method (ast, regex)")
	totalWorkers := flag.Int("workers", 4, "Number of workers")
	timeout := flag.Int("timeout", 10, "Timeout in seconds")
	patternsFilePath := flag.String("patterns", "", "Path to file with PII patterns")
	outputFilePath := flag.String("output", "", "Path to output file; default: stdout; supported formats: *.json")
	flag.Parse()

	patterns := make([]string, 0)
	if *patternsFilePath != "" {
		parsedPatterns, err := config.ReadPIIPatterns(*patternsFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		patterns = parsedPatterns
	}

	config.Setup(
		config.WithMethod(piidetect.DetectMethod(*method)),
		config.WithTotalWorkers(*totalWorkers),
		config.WithTimeout(time.Duration(*timeout)*time.Second),
		config.WithPIIPatterns(patterns),
		config.WithOutputFilePath(*outputFilePath),
	)

	// TODO: accept paths as arguments
	rawPaths := make([]string, 0)
	if len(rawPaths) == 0 {
		rawPaths = append(rawPaths, "./...")
	}

	paths, err := getExpandedPaths(rawPaths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	issues, err := getIssues(paths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(issues) == 0 {
		fmt.Println("No PII issues found.")
		return
	}

	report.Generate(issues)
}

func getExpandedPaths(rawPaths []string) ([]string, error) {
	paths := make([]string, 0)

	for _, rawPath := range rawPaths {
		if rawPath == "./..." {
			rawPath = "."
		}
		expandedPaths := make([]string, 0)
		if err := filepath.Walk(rawPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && !strings.Contains(path, "_test.go") {
				paths = append(paths, path)
			}
			return nil

		}); err != nil {
			return nil, err
		}
		paths = append(paths, expandedPaths...)
	}

	return paths, nil
}

func getIssues(paths []string) ([]piidetect.Issue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(config.TotalWorkers())

	inputChan := make(chan string, len(paths))
	outputChan := make(chan piidetect.Issue, config.TotalWorkers())

	for i := 0; i < config.TotalWorkers(); i++ {
		go func() {
			defer wg.Done()
			for path := range inputChan {
				newIssues, err := detector.DetectPII(ctx, path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "(%s) Error: %v\n", path, err)
					continue
				}
				for _, issue := range newIssues {
					outputChan <- issue
				}
			}
		}()
	}

	for _, path := range paths {
		inputChan <- path
	}
	close(inputChan)

	go func() {
		wg.Wait()
		close(outputChan)
	}()

	issues := make([]piidetect.Issue, 0)
	for issue := range outputChan {
		issues = append(issues, issue)
	}
	return issues, nil
}
