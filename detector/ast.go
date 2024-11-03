package detector

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/ullauri/piidetect"
)

func isGoFile(path string) bool {
	return len(path) >= 3 && path[len(path)-3:] == ".go"
}

func detectPIIAST(_ context.Context, path string) ([]piidetect.Issue, error) {
	if !isGoFile(path) {
		return nil, nil
	}

	issues := make([]piidetect.Issue, 0)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			if issue := checkAST(callExpr, fset); issue != nil {
				issues = append(issues, *issue)
			}
		}
		return true
	})

	return issues, nil
}

func checkAST(callExpr *ast.CallExpr, fset *token.FileSet) *piidetect.Issue {
	for _, arg := range callExpr.Args {
		switch expr := arg.(type) {
		case *ast.Ident:
			if match := containsPII(expr.Name); match != nil {
				pos := fset.Position(expr.Pos())
				return &piidetect.Issue{
					Match:   *match,
					Type:    piidetect.Identifier,
					File:    pos.Filename,
					Line:    pos.Line,
					Message: fmt.Sprintf(expr.Name),
				}
			}
		case *ast.BasicLit:
			if expr.Kind == token.STRING {
				if match := containsPII(expr.Value); match != nil {
					pos := fset.Position(expr.Pos())
					return &piidetect.Issue{
						Match:   *match,
						Type:    piidetect.LiteralString,
						File:    pos.Filename,
						Line:    pos.Line,
						Message: fmt.Sprintf(expr.Value),
					}
				}
			}
		case *ast.SelectorExpr:
			if ident, ok := expr.X.(*ast.Ident); ok {
				if match := containsPII(ident.Name + "." + expr.Sel.Name); match != nil {
					pos := fset.Position(expr.Pos())
					return &piidetect.Issue{
						Match:   *match,
						Type:    piidetect.StructField,
						File:    pos.Filename,
						Line:    pos.Line,
						Message: fmt.Sprintf("%s.%s", ident.Name, expr.Sel.Name),
					}
				}
			}
		}
	}
	return nil
}
