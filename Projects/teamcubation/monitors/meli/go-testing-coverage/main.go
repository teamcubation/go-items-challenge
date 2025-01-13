package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	CommitID string  `json:"commit_id"`
	File     string  `json:"file"`
	Line     int     `json:"line"`
	Coverage float64 `json:"coverage"`
}

func getModuleName(repoPath string) (string, error) {
	goModPath := filepath.Join(repoPath, "go.mod")
	file, err := os.Open(goModPath)
	if err != nil {
		return "", fmt.Errorf("error opening go.mod: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				return parts[1], nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading go.mod: %v", err)
	}

	return "", fmt.Errorf("module name not found in go.mod")
}

func readCovignore(repoPath string) (map[string]bool, error) {
	covignoreFilePath := filepath.Join(repoPath, ".covignore")
	covignoreData, err := os.ReadFile(covignoreFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading covignore file: %v", err)
	}

	ignoreFiles := make(map[string]bool)
	lines := strings.Split(string(covignoreData), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ignoreFiles[line] = true
	}

	return ignoreFiles, nil
}

func shouldIgnore(file, moduleName string, ignoreFiles map[string]bool) bool {
	fileWithoutModule := strings.TrimPrefix(file, fmt.Sprintf("%s/", moduleName))

	if ignoreFiles[fileWithoutModule] {
		return true
	}

	for ignorePattern := range ignoreFiles {
		match, _ := filepath.Match(ignorePattern, fileWithoutModule)
		if match {
			return true
		}
	}
	return false
}

func analyzeCoverage(repoPath string) (Metric, error) {
	metric := Metric{
		MetricID: "coverage",
		Score:    "1",
		Evidence: []Evidence{},
	}

	// Run go test with coverage
	cmd := exec.Command("go", "test", "./...", "-coverprofile=coverage.out")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		return metric, fmt.Errorf("error running tests: %v", err)
	}

	ignoreFiles, err := readCovignore(repoPath)
	if err != nil {
		log.Printf("Error reading covignore: %v", err)
		ignoreFiles = make(map[string]bool)
	}

	moduleName, err := getModuleName(repoPath)
	if err != nil {
		log.Printf("Error getting module name: %v", err)
		return metric, err
	}

	// Parse coverage output
	coverageFilePath := filepath.Join(repoPath, "coverage.out")
	coverageData, err := os.ReadFile(coverageFilePath)
	if err != nil {
		return metric, fmt.Errorf("error reading coverage file: %v", err)
	}

	lines := strings.Split(string(coverageData), "\n")
	if len(lines) == 0 {
		return metric, fmt.Errorf("coverage file is empty")
	}

	filesCoverage := make(map[string]struct {
		statements int
		covered    int
	})

	for _, line := range lines {
		if strings.HasPrefix(line, "mode:") || line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			continue
		}

		fileParts := strings.Split(parts[0], ":")
		if len(fileParts) < 2 {
			continue
		}

		file := fileParts[0]
		if shouldIgnore(file, moduleName, ignoreFiles) {
			continue
		}

		statements, _ := strconv.Atoi(parts[1])
		executions, _ := strconv.Atoi(parts[2])

		coverage := filesCoverage[file]
		coverage.statements += statements
		if executions > 0 {
			coverage.covered += statements
		}
		filesCoverage[file] = coverage
	}

	totalStatements := 0
	totalCovered := 0

	for file, coverage := range filesCoverage {
		totalStatements += coverage.statements
		totalCovered += coverage.covered

		fileCoverage := float64(coverage.covered) / float64(coverage.statements) * 100

		author, commitID, err := getFileInfo(repoPath, file)
		if err != nil {
			log.Printf("Error getting file info for %s: %v", file, err)
		}

		metric.Evidence = append(metric.Evidence, Evidence{
			CommitID: commitID,
			File:     file,
			Line:     0,
			Coverage: fileCoverage,
		})

		if metric.GitAuthor == "" {
			metric.GitAuthor = author
		}
	}

	if totalStatements > 0 {
		overallCoverage := float64(totalCovered) / float64(totalStatements) * 100
		if overallCoverage >= 80 {
			metric.Score = "3"
		} else if overallCoverage >= 50 {
			metric.Score = "2"
		}
	}

	return metric, nil
}

func getFileInfo(repoPath, file string) (string, string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", "", fmt.Errorf("error opening repository: %v", err)
	}

	commits, err := repo.Log(&git.LogOptions{FileName: &file})
	if err != nil {
		return "", "", fmt.Errorf("error getting file log: %v", err)
	}

	commit, err := commits.Next()
	if err != nil {
		return "", "", fmt.Errorf("error getting commit: %v", err)
	}

	return commit.Author.Email, commit.Hash.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <repo_path>")
		return
	}

	repoPath := os.Args[1]

	coverageMetric, err := analyzeCoverage(repoPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	jsonOutput, err := json.MarshalIndent([]Metric{coverageMetric}, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	fmt.Println(string(jsonOutput))
}
