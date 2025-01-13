package analysis

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func AnalyzeErrorWrapping(fileContent, filePath string) FileAnalysisResult {
	var evidence []Evidence
	if strings.Contains(filePath, "cmd") {
		return FileAnalysisResult{IsValid: false}
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", fileContent, parser.AllErrors)
	if err != nil {
		fmt.Println("Failed to parse source:", err)
		return FileAnalysisResult{IsValid: false}
	}

	totalErrors := 0
	wrappedErrors := 0

	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if retStmt, ok := n.(*ast.ReturnStmt); ok {
					for _, result := range retStmt.Results {
						if ident, ok := result.(*ast.Ident); ok && ident.Name == "err" {
							totalErrors++

							if isErrorWrapped(fn.Body, ident) {
								wrappedErrors++
							} else {
								evidence = append(evidence, Evidence{
									File: filePath,
									Line: fset.Position(retStmt.Pos()).Line,
								})
							}
						}
					}
				}
				return true
			})
		}
		return true
	})

	fileAnalysisResult := calculateScore(totalErrors, wrappedErrors)
	fileAnalysisResult.Evidences = evidence
	return fileAnalysisResult
}

func isErrorWrapped(body *ast.BlockStmt, errIdent *ast.Ident) bool {
	isWrapped := false

	ast.Inspect(body, func(n ast.Node) bool {
		// Verificar si el error es envuelto.
		if retStmt, ok := n.(*ast.ReturnStmt); ok {
			for _, result := range retStmt.Results {
				if callExpr, ok := result.(*ast.CallExpr); ok {
					// Verificar si se está utilizando fmt.Errorf o errors.Wrap.
					if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
						if ident, ok := fun.X.(*ast.Ident); ok && ident.Name == "fmt" && (fun.Sel.Name == "Errorf" || fun.Sel.Name == "Wrap") {
							// Verificar que el error original (errIdent) sea parte del envoltorio.
							for _, arg := range callExpr.Args {
								if id, ok := arg.(*ast.Ident); ok && id.Name == errIdent.Name {
									isWrapped = true
									return false
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	return isWrapped
}

func calculateScore(totalErrors, wrappedErrors int) FileAnalysisResult {
	if totalErrors == 0 {
		return FileAnalysisResult{IsValid: false, Score: 0}
	}
	wrapRatio := float64(wrappedErrors) / float64(totalErrors)

	switch {
	case wrapRatio == 1:
		return FileAnalysisResult{IsValid: true, Score: 5}
	case wrapRatio >= 0.75:
		return FileAnalysisResult{IsValid: true, Score: 4}
	case wrapRatio >= 0.5:
		return FileAnalysisResult{IsValid: true, Score: 3}
	case wrapRatio >= 0.25:
		return FileAnalysisResult{IsValid: true, Score: 2}
	default:
		return FileAnalysisResult{IsValid: true, Score: 1}
	}
}
