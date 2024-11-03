package detector

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/ullauri/piidetect"
	"github.com/ullauri/piidetect/config"
)

func detectPIIRegex(ctx context.Context, path string) ([]piidetect.Issue, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var issues []piidetect.Issue
	scanner := bufio.NewScanner(file)
	lineNumber := 1

	var patterns []*regexp.Regexp
	for _, keyword := range config.PIIPatterns() {
		pattern, err := regexp.Compile("(?i)\\b" + regexp.QuoteMeta(keyword) + "\\b")
		if err != nil {
			return nil, fmt.Errorf("error compiling regex for pattern %q: %v", keyword, err)
		}
		patterns = append(patterns, pattern)
	}

	for scanner.Scan() {
		line := scanner.Text()

		for _, pattern := range patterns {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				if pattern.MatchString(line) {
					issues = append(issues, piidetect.Issue{
						Match:   pattern.String(),
						Type:    piidetect.LiteralString,
						File:    path,
						Line:    lineNumber,
						Message: fmt.Sprintf("%q", line),
					})
				}
			}
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return issues, nil
}
