package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Evidence struct {
	Error string `json:"error"`
}

type Metric struct {
	MetricID string     `json:"metric_id"`
	Score    string     `json:"score"`
	Evidence []Evidence `json:"evidence,omitempty"`
}

func checkDirExists(basePath, dirName string) bool {
	path := filepath.Join(basePath, dirName)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <repo_path> [modified_files...]")
		os.Exit(1)
	}

	repoPath := os.Args[1]

	metric := Metric{
		MetricID: "initial-scaffolding",
		Score:    "3",
		Evidence: []Evidence{},
	}

	score := 3
	var evidences []Evidence

	mainPath := filepath.Join("cmd", "api", "main.go")
	if !checkDirExists(repoPath, mainPath) {
		score = 1
		evidences = append(evidences, Evidence{Error: "main.go file does not exist in cmd/api"})
	}

	folders := []string{"internal", "pkg"}
	exist := false

	for _, folder := range folders {
		if checkDirExists(repoPath, folder) {
			exist = true
			break
		}
	}

	if !exist {
		score = 2
		evidences = append(evidences, Evidence{Error: "neither internal nor pkg directory exists"})
	}

	metric.Score = fmt.Sprintf("%d", score)
	if len(evidences) > 0 {
		metric.Evidence = evidences
	}

	output, _ := json.MarshalIndent([]Metric{metric}, "", "  ")
	fmt.Println(string(output))
}
