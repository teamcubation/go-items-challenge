package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

const (
	melisourceRepo   string = "github.com/melisource"
	mercadolibreRepo string = "github.com/mercadolibre"
)

var (
	msRustyImport      = fmt.Sprintf("%s/fury_go-core/pkg/rusty", melisourceRepo)
	mlRustyImport      = fmt.Sprintf("%s/fury_go-core/pkg/rusty", mercadolibreRepo)
	msHttpClientImport = fmt.Sprintf("%s/fury_go-core/pkg/transport/httpclient", melisourceRepo)
	mlHttpClientImport = fmt.Sprintf("%s/fury_go-core/pkg/transport/httpclient", mercadolibreRepo)
	msBreakerImport    = fmt.Sprintf("%s/fury_go-core/pkg/breaker", melisourceRepo)
	mlBreakerImport    = fmt.Sprintf("%s/fury_go-core/pkg/breaker", mercadolibreRepo)
)

type Metric struct {
	MetricID  string     `json:"metric_id"`
	GitAuthor string     `json:"git_author"`
	Score     string     `json:"score"`
	Evidence  []Evidence `json:"evidence"`
}

type Evidence struct {
	CommitID string `json:"commit_id"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type Skill struct {
	ID   string
	Name string
}

type SkillData struct {
	Score    int
	Evidence []Evidence
}

type Results map[string]map[string]SkillData

var skills = []Skill{
	{ID: "rusty_pkg_usage", Name: "Proper usage of rusty for making API requests"},
	{ID: "http_client_provided", Name: "httpclient pkg Usage"},
	{ID: "rusty_error_handling", Name: "Handling errors effectively in REST client requests"},
	{ID: "http_retry_mechanism", Name: "Implementing retry mechanisms for HTTP requests"},
	{ID: "circuit_breaker_pattern", Name: "Using circuit breaker pattern for REST requests"},
}

func analyzeRepo(repoPath string, files []string) ([]Metric, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		fmt.Printf("Error opening repository: %v\n", err)
		return nil, err
	}

	filesToAnalyze, err := getFilesToAnalyze(repo, files)
	if err != nil {
		return nil, err
	}

	results := make(map[string]map[string]SkillData)

	for _, filePath := range filesToAnalyze {
		fullPath := filepath.Join(repoPath, filePath)
		author, err := getFileAuthor(repo, filePath)
		if err != nil {
			return nil, err
		}

		commitID, err := getCommitID(repo, filePath)
		if err != nil {
			return nil, err
		}

		fileResults, err := analyzeFile(fullPath)
		if err != nil {
			return nil, err
		}

		if _, exists := results[author]; !exists {
			results[author] = make(map[string]SkillData)
			for _, skill := range skills {
				results[author][skill.ID] = SkillData{Score: 0, Evidence: []Evidence{}}
			}
		}

		for skillID, data := range fileResults {
			skillData := results[author][skillID]

			if data.Score > skillData.Score {
				skillData.Score = data.Score
				skillData.Evidence = []Evidence{}
			}

			if data.Score == skillData.Score {
				for _, line := range data.Evidence {
					skillData.Evidence = append(skillData.Evidence, Evidence{
						CommitID: commitID,
						File:     filePath,
						Line:     line.Line,
					})
				}
			}

			results[author][skillID] = skillData
		}
	}

	var output []Metric
	for author, scores := range results {
		for _, skill := range skills {
			skillID := skill.ID
			output = append(output, Metric{
				MetricID:  skillID,
				GitAuthor: author,
				Score:     fmt.Sprintf("%d", scores[skillID].Score),
				Evidence:  scores[skillID].Evidence,
			})
		}
	}

	return output, nil
}

func analyzeFile(filePath string) (map[string]SkillData, error) {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	results := make(map[string]SkillData)
	for _, skill := range skills {
		results[skill.ID] = SkillData{Score: 0, Evidence: []Evidence{}}
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			if strings.Contains(x.Path.Value, msRustyImport) || strings.Contains(x.Path.Value, mlRustyImport) {
				markSkill(results, "rusty_pkg_usage", fset.Position(n.Pos()).Line)
			}
			if strings.Contains(x.Path.Value, msHttpClientImport) || strings.Contains(x.Path.Value, mlHttpClientImport) {
				markSkill(results, "http_client_provided", fset.Position(n.Pos()).Line)
			}
			if strings.Contains(x.Path.Value, msBreakerImport) || strings.Contains(x.Path.Value, mlBreakerImport) {
				markSkill(results, "circuit_breaker_pattern", fset.Position(n.Pos()).Line)
			}
		case *ast.CallExpr:
			if fun, ok := x.Fun.(*ast.SelectorExpr); ok {
				if pkgIdent, ok := fun.X.(*ast.Ident); ok && pkgIdent.Name == "httpclient" && fun.Sel.Name == "NewRetryable" {
					markSkill(results, "http_retry_mechanism", fset.Position(n.Pos()).Line)
				}
			}
		case *ast.IfStmt:
			if binExpr, ok := x.Cond.(*ast.BinaryExpr); ok {
				if binExpr.Op == token.NEQ && isNilComparison(binExpr) {
					foundRustyError := false
					foundErrorsAs := false

					ast.Inspect(x.Body, func(n ast.Node) bool {
						switch stmt := n.(type) {
						case *ast.DeclStmt:
							if genDecl, ok := stmt.Decl.(*ast.GenDecl); ok {
								for _, spec := range genDecl.Specs {
									if valueSpec, ok := spec.(*ast.ValueSpec); ok {
										if starExpr, ok := valueSpec.Type.(*ast.StarExpr); ok {
											if selExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
												if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "rusty" && selExpr.Sel.Name == "Error" {
													foundRustyError = true
												}
											}
										}
									}
								}
							}
						case *ast.IfStmt:
							if callExpr, ok := stmt.Cond.(*ast.CallExpr); ok {
								if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
									if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "errors" && selExpr.Sel.Name == "As" {
										foundErrorsAs = true
									}
								}
							}
						}
						return true
					})
					if foundRustyError && foundErrorsAs {
						markSkill(results, "rusty_error_handling", fset.Position(x.Pos()).Line)
					}
				}
			}
		}
		return true
	})

	return results, nil
}

func isNilComparison(expr ast.Expr) bool {
	if binExpr, ok := expr.(*ast.BinaryExpr); ok {
		return (isNilIdent(binExpr.X) && !isNilIdent(binExpr.Y)) ||
			(!isNilIdent(binExpr.X) && isNilIdent(binExpr.Y))
	}
	return false
}

func isNilIdent(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == "nil"
	}
	return false
}

func markSkill(results map[string]SkillData, skillID string, line int) {
	skillData := results[skillID]
	skillData.Score = 1
	skillData.Evidence = append(skillData.Evidence, Evidence{Line: line})
	results[skillID] = skillData
}

func getFilesToAnalyze(repo *git.Repository, files []string) ([]string, error) {
	var filesToAnalyze []string

	if len(files) == 0 {
		worktree, err := repo.Worktree()
		if err != nil {
			return nil, err
		}

		err = filepath.Walk(worktree.Filesystem.Root(), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".go" {
				relPath, err := filepath.Rel(worktree.Filesystem.Root(), path)
				if err != nil {
					return err
				}
				filesToAnalyze = append(filesToAnalyze, relPath)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		for _, file := range files {
			if strings.HasSuffix(file, ".go") {
				filesToAnalyze = append(filesToAnalyze, file)
			}
		}
	}

	return filesToAnalyze, nil
}

func getFileAuthor(repo *git.Repository, file string) (string, error) {
	commits, err := repo.Log(&git.LogOptions{FileName: &file})
	if err != nil {
		return "", err
	}

	commit, err := commits.Next()
	if err != nil {
		return "", err
	}

	return commit.Author.Email, nil
}

func getCommitID(repo *git.Repository, file string) (string, error) {
	commits, err := repo.Log(&git.LogOptions{FileName: &file})
	if err != nil {
		return "", err
	}

	objectCommit, err := commits.Next()
	if err != nil {
		return "", err
	}

	return objectCommit.Hash.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <repo_path> [file1] [file2] ...")
		return
	}

	repoPath := os.Args[1]
	filesToAnalyze := os.Args[2:]

	metrics, err := analyzeRepo(repoPath, filesToAnalyze)
	if err != nil {
		log.Fatal(err)
		return
	}

	jsonOutput, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	fmt.Println(string(jsonOutput))
}
