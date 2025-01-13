package main

import (
	"encoding/json"
	"errhandlmon/internal/monitor"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <repo_path> [modified_files...]")
		os.Exit(1)
	}

	repoPath := os.Args[1]

	var modifiedFiles []string
	if len(os.Args) > 2 {
		modifiedFiles = os.Args[1:]
	}

	results, err := monitor.AnalyzeRepo(repoPath, modifiedFiles)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	jsonResults, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonResults))
}
