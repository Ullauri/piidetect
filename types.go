package piidetect

type DetectMethod string

const (
	AST   DetectMethod = "ast"
	Regex DetectMethod = "regex"
)

type IssueType string

const (
	Identifier    IssueType = "identifier"
	LiteralString IssueType = "literal_string"
	StructField   IssueType = "struct_field"
)

type Issue struct {
	Match   string
	Type    IssueType
	File    string
	Line    int
	Message string
}
