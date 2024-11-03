package config

import (
	"sync"
	"time"

	"github.com/ullauri/piidetect"
)

type config struct {
	Timeout        time.Duration
	TotalWorkers   int
	Method         piidetect.DetectMethod
	PIIPatterns    []string
	OutputFilePath string
}

var internalConfig *config
var configMutex sync.Mutex

func DefaultSetup() {
	Setup(
		WithPIIPatterns([]string{}),
		WithMethod(""),
		WithTotalWorkers(0),
		WithTimeout(0),
		WithOutputFilePath(""),
	)
}

func Setup(opts ...func(*config)) {
	configMutex.Lock()
	defer configMutex.Unlock()

	newConfig := &config{}
	for _, opt := range opts {
		opt(newConfig)
	}

	internalConfig = newConfig
}

func WithPIIPatterns(patterns []string) func(*config) {
	return func(c *config) {
		if len(patterns) == 0 {
			patterns = defaultPatterns
		}
		c.PIIPatterns = patterns
	}
}

func WithMethod(method piidetect.DetectMethod) func(*config) {
	return func(c *config) {
		if method == "" {
			method = piidetect.AST
		}
		c.Method = method
	}
}

func WithTotalWorkers(workers int) func(*config) {
	return func(c *config) {
		if workers <= 0 {
			workers = 10
		}
		c.TotalWorkers = workers
	}
}

func WithTimeout(timeout time.Duration) func(*config) {
	return func(c *config) {
		if timeout <= 0 {
			timeout = 10 * time.Second
		}
		c.Timeout = timeout
	}
}

func WithOutputFilePath(path string) func(*config) {
	return func(c *config) {
		c.OutputFilePath = path
	}
}

func PIIPatterns() []string {
	return internalConfig.PIIPatterns
}

func Method() piidetect.DetectMethod {
	return internalConfig.Method
}

func TotalWorkers() int {
	return internalConfig.TotalWorkers
}

func Timeout() time.Duration {
	return internalConfig.Timeout
}

func OutputFilePath() string {
	return internalConfig.OutputFilePath
}
