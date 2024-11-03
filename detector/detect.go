package detector

import (
	"context"
	"fmt"
	"strings"

	"github.com/ullauri/piidetect"
	"github.com/ullauri/piidetect/config"
)

func DetectPII(ctx context.Context, path string) ([]piidetect.Issue, error) {
	switch config.Method() {
	case piidetect.AST:
		return detectPIIAST(ctx, path)
	case piidetect.Regex:
		return detectPIIRegex(ctx, path)
	default:
		panic(fmt.Sprintf("unsupported detection method: %s", config.Method()))
	}
}

func containsPII(value string) *string {
	value = strings.Trim(value, "\"'")
	for _, keyword := range config.PIIPatterns() {
		if strings.Contains(strings.ToLower(value), keyword) {
			return &keyword
		}
	}
	return nil
}
