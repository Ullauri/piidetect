# PII Detect

`piidetect` is a Go-based static analysis tool designed to identify instances of Personally Identifiable Information (PII) within source code files. It leverages either Abstract Syntax Tree (AST) or regular expression (regex) methods to detect sensitive data patterns. This tool supports custom patterns, configurable worker concurrency, and includes options for reporting detected issues.

## Features

- **Customizable Detection Methods**: Supports both `ast` and `regex` detection methods.
- **Configurable Worker Pool**: Set the number of concurrent workers to optimize analysis.
- **Flexible Pattern Configuration**: Load PII patterns from a custom file or use built-in patterns.
- **Detailed Reporting**: Generates a report of PII instances detected within the provided file paths.

## Installation

To install `piidetect`, use the `go install` command:

```bash
go install github.com/ullauri/piidetect/cmd/piidetect@latest
```

This will install piidetect as a command-line tool, which you can then run from any directory.

## Usage

Run piidetect by specifying your detection method, number of workers, timeout, and PII patterns file.

```bash
piidetect ./...
```

## Command-Line Flags
-method: Specifies the detection method (ast or regex). Default is ast.
- workers: Sets the number of concurrent worker goroutines. Default is 4.
- timeout: Sets a timeout for the analysis (in seconds). Default is 10 seconds.
- patterns: Path to a file containing custom PII patterns, each pattern on a new line.
- output: Optional path to output file.

