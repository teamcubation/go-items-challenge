package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
)

var (
	restClientImportRegex   = regexp.MustCompile(`import\s+.*com\.mercadolibre\.restclient\.MeliRestClient`)
	restPoolImportRegex     = regexp.MustCompile(`import\s+.*com\.mercadolibre\.restclient\.RESTPool`)
	restPoolCreationRegex   = regexp.MustCompile(`RESTPool\.builder\(\)`)
	restClientCreationRegex = regexp.MustCompile(`MeliRestClient\.builder\(\)`)
	errorHandlingRegex      = regexp.MustCompile(`catch\s*\((RestException|ParseException)\s*(\w+)\)`)
	retryStrategyRegex      = regexp.MustCompile(`withRetryStrategy\([^)]*\)`)
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
	{ID: "rest_client_usage", Name: "REST Client Usage"},
	{ID: "rest_pool_usage", Name: "REST Pool Usage"},
	{ID: "error_handling", Name: "Handling exceptions effectively in REST client requests"},
	{ID: "retry_strategy_usage", Name: "Implementing retry mechanisms for HTTP requests"},
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
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	results := make(map[string]SkillData)
	for _, skill := range skills {
		results[skill.ID] = SkillData{Score: 0, Evidence: []Evidence{}}
	}

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		if restClientImportRegex.MatchString(line) {
			markSkill(results, "rest_client_usage", lineNumber)
		}

		if restPoolImportRegex.MatchString(line) {
			markSkill(results, "rest_pool_usage", lineNumber)
		}

		if restPoolCreationRegex.MatchString(line) {
			markSkill(results, "rest_pool_usage", lineNumber)
		}

		if restClientCreationRegex.MatchString(line) {
			markSkill(results, "rest_client_usage", lineNumber)
		}

		if errorHandlingRegex.MatchString(line) {
			markSkill(results, "error_handling", lineNumber)
		}

		if retryStrategyRegex.MatchString(line) {
			markSkill(results, "retry_strategy_usage", lineNumber)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
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
			if !info.IsDir() && filepath.Ext(path) == ".java" {
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
