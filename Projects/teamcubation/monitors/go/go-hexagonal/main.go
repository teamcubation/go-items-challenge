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

type InterfaceSpec struct {
	Methods map[string]struct{}
}

func getInterfacesFromFile(filePath string) (map[string]InterfaceSpec, error) {
	interfaces := make(map[string]InterfaceSpec)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if t, ok := n.(*ast.TypeSpec); ok {
			if interfaceType, ok := t.Type.(*ast.InterfaceType); ok {
				interfaceName := t.Name.Name
				methods := make(map[string]struct{})

				for _, method := range interfaceType.Methods.List {
					for _, name := range method.Names {
						methods[name.Name] = struct{}{}
					}
				}

				interfaces[interfaceName] = InterfaceSpec{
					Methods: methods,
				}
			}
		}
		return true
	})

	return interfaces, nil
}

func checkImplementationsInAdapters(repoPath string, interfaces map[string]InterfaceSpec) (map[string]Metric, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	results := make(map[string]Metric)
	implementedInterfaces := make(map[string]map[string]bool)
	adaptersPath := filepath.Join(repoPath, "internal/adapters")

	err = filepath.Walk(adaptersPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(repoPath, path)
			if err != nil {
				return err
			}

			author, err := getFileAuthor(repo, relPath)
			if err != nil {
				return err
			}

			commitID, err := getCommitID(repo, relPath)
			if err != nil {
				return err
			}

			if _, exists := implementedInterfaces[author]; !exists {
				implementedInterfaces[author] = make(map[string]bool)
			}

			ast.Inspect(file, func(n ast.Node) bool {
				typeSpec, ok := n.(*ast.TypeSpec)
				if !ok {
					return true
				}

				_, ok = typeSpec.Type.(*ast.StructType)
				if !ok {
					return true
				}

				structName := typeSpec.Name.Name
				structMethods := make(map[string]struct{})

				// Collect methods for this struct
				ast.Inspect(file, func(n ast.Node) bool {
					funcDecl, ok := n.(*ast.FuncDecl)
					if !ok || funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
						return true
					}

					starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr)
					if !ok {
						return true
					}

					ident, ok := starExpr.X.(*ast.Ident)
					if !ok || ident.Name != structName {
						return true
					}

					structMethods[funcDecl.Name.Name] = struct{}{}
					return true
				})

				// Check if this struct implements any of the interfaces
				for ifaceName, ifaceSpec := range interfaces {
					allMethodsImplemented := true
					for method := range ifaceSpec.Methods {
						if _, ok := structMethods[method]; !ok {
							allMethodsImplemented = false
							break
						}
					}

					if allMethodsImplemented {
						implementedInterfaces[author][ifaceName] = true
						metric, exists := results[author]
						if !exists {
							metric = Metric{
								MetricID:  "port_implementation",
								GitAuthor: author,
								Evidence:  []Evidence{},
							}
						}

						metric.Evidence = append(metric.Evidence, Evidence{
							CommitID: commitID,
							File:     relPath,
							Line:     fset.Position(typeSpec.Pos()).Line,
						})
						results[author] = metric
					}
				}

				return true
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Calculate scores
	for author, implementedIfaces := range implementedInterfaces {
		metric := results[author]
		implementedCount := len(implementedIfaces)
		totalInterfaces := len(interfaces)

		switch {
		case implementedCount == totalInterfaces:
			metric.Score = "3"
		case implementedCount > 0:
			metric.Score = "2"
		default:
			metric.Score = "1"
		}

		results[author] = metric
	}

	return results, nil
}

func analyzeRepo(repoPath string, modifiedFiles []string) ([]Metric, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		fmt.Printf("Error opening repository: %v\n", err)
		return nil, err
	}

	hexagonalMetric := analyzeHexagonalScaffolding(repoPath)

	filesToAnalyze, err := getFilesToAnalyze(repo, modifiedFiles)
	if err != nil {
		return nil, err
	}

	interfaces := make(map[string]InterfaceSpec)

	for _, filePath := range filesToAnalyze {
		fullPath := filepath.Join(repoPath, filePath)

		if strings.Contains(filePath, "internal/core/ports") {
			fileInterfaces, err := getInterfacesFromFile(fullPath)
			if err != nil {
				fmt.Printf("Error analyzing file %s: %v\n", fullPath, err)
				continue
			}
			interfaces = fileInterfaces
		}
	}

	var output []Metric

	if len(interfaces) == 0 {
		output = append(output, hexagonalMetric, Metric{
			MetricID: "port_implementation",
			Score:    "N/A",
		})

		return output, nil
	}

	results, err := checkImplementationsInAdapters(repoPath, interfaces)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}

	for _, metric := range results {
		output = append(output, metric)
	}

	output = append(output, hexagonalMetric)

	return output, nil
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

func directoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
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

func analyzeHexagonalScaffolding(repoPath string) Metric {
	requiredPackages := []string{"internal/core/ports", "internal/core/domain", "internal/adapters"}
	existingPackages := 0

	metric := Metric{
		MetricID: "hexagonal_scaffolding",
		Score:    "3",
	}

	for _, pkg := range requiredPackages {
		fullPath := filepath.Join(repoPath, pkg)
		exists, err := directoryExists(fullPath)
		if err != nil || !exists {
			fmt.Printf("Error checking directory %s: %v\n", fullPath, err)

			metric.Evidence = append(metric.Evidence, Evidence{
				File: pkg,
				Line: 0,
			})

			continue
		}

		existingPackages++
	}

	if existingPackages == 0 {
		metric.Score = "1"
	} else if existingPackages < len(requiredPackages) {
		metric.Score = "2"
	}

	return metric
}
