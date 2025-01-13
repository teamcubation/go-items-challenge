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
	{ID: "go_unit_tests", Name: "Writing unit tests using Go's testing package"},
	{ID: "testify_assertions", Name: "Using Testify for enhanced assertions"},
	{ID: "table_driven_tests", Name: "Implementing table-driven tests"},
	{ID: "mock_objects", Name: "Using mock objects for testing"},
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
		if fileInfo, err := os.Stat(fullPath); err == nil && !fileInfo.IsDir() && filepath.Ext(filePath) == ".go" {
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
	}

	var output []Metric
	for author, scores := range results {
		for _, skill := range skills {
			skillID := skill.ID
			output = append(output, Metric{
				MetricID:  skillID,
				GitAuthor: author,
				Score:     string(scores[skillID].Score),
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
	testFuncRegex := regexp.MustCompile(`^func Test\w+\(t \*testing\.T\)`)
	testifyRegex := regexp.MustCompile(`github.com/stretchr/testify`)
	tableDrivenTestSliceRegex := regexp.MustCompile(`^\s*test\s*:=\s*\[\s*struct\s*{\s*name\s+string`)
	tableDrivenTestLoopRegex := regexp.MustCompile(`^\s*for\s+_,\s+tt\s+:=\s+range\s+test\s*{`)

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		for _, skill := range skills {
			skillData := results[skill.ID]
			switch skill.ID {
			case "go_unit_tests":
				if testFuncRegex.MatchString(line) {
					skillData.Score = 1
					skillData.Evidence = append(skillData.Evidence, Evidence{Line: lineNumber})
				}
			case "testify_assertions":
				if testifyRegex.MatchString(line) {
					skillData.Score = 1
					skillData.Evidence = append(skillData.Evidence, Evidence{Line: lineNumber})
				}
			case "mock_objects":
				if strings.Contains(line, "github.com/golang/mock/gomock") ||
					strings.Contains(line, "go.uber.org/mock/gomock") ||
					strings.Contains(line, "github.com/stretchr/testify/mock") {
					skillData.Score = 1
					skillData.Evidence = append(skillData.Evidence, Evidence{Line: lineNumber})
				}
			case "table_driven_tests":
				if tableDrivenTestSliceRegex.MatchString(line) {
					skillData.Score = 1
					skillData.Evidence = append(skillData.Evidence, Evidence{Line: lineNumber})
				}
				if tableDrivenTestLoopRegex.MatchString(line) {
					skillData.Score = 1
					skillData.Evidence = append(skillData.Evidence, Evidence{Line: lineNumber})
				}
			}
			results[skill.ID] = skillData
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
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
			if !info.IsDir() && filepath.Ext(path) == ".go" && strings.HasSuffix(path, "_test.go") {
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
			if strings.HasSuffix(file, "_test.go") {
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

	return fmt.Sprintf("%s <%s>", commit.Author.Name, commit.Author.Email), nil
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
