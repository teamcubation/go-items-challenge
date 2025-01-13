package monitor

import (
	"errhandlmon/internal/monitor/analysis"
	"errhandlmon/internal/monitor/file"
	"errhandlmon/internal/monitor/git"
	"fmt"
	"os"
	"path/filepath"
)

type Result struct {
	MetricID  string              `json:"metric_id"`
	GitAuthor string              `json:"git_author"`
	Score     string              `json:"score"`
	Evidence  []analysis.Evidence `json:"evidence"`
}

func AnalyzeRepo(repoPath string, modifiedFiles []string) (map[string][]Result, error) {
	results := make(map[string][]Result)

	if len(modifiedFiles) == 0 {
		var err error
		modifiedFiles, err = file.GetAllGoFiles(repoPath)
		if err != nil {
			return nil, fmt.Errorf("AnalyzeRepo: failed getting all files from repo %s: %w", repoPath, err)
		}
	}

	var score int

	for _, filePath := range modifiedFiles {
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("AnalyzeRepo: error reading file %s: %w", filePath, err)
		}

		relativePath, err := getRelativePath(repoPath, filePath)
		if err != nil {
			return nil, fmt.Errorf("AnalyzeRepo: error getting relative path: %w", err)
		}

		panicEvidence := analysis.AnalyzePanicUsage(string(fileContent), relativePath)
		wrapEvidence := analysis.AnalyzeErrorWrapping(string(fileContent), relativePath)
		ignoreEvidence := analysis.AnalyzeErrorsIgnore(string(fileContent), relativePath)

		if len(panicEvidence) > 0 || len(wrapEvidence.Evidences) > 0 || len(ignoreEvidence) > 0 {
			commitHash, err := git.GetLastCommitHash(relativePath, repoPath)
			if err != nil {
				return nil, fmt.Errorf("AnalyzeRepo: error getting last commit from file %s: %w", filePath, err)
			}

			author, err := git.GetGitAuthor(commitHash, repoPath)
			if err != nil {
				return nil, fmt.Errorf("AnalyzeRepo: error getting git author from file %s: %w", filePath, err)
			}

			addPanicUsageResults(results, author, panicEvidence)
			addErrorWrapResults(results, author, wrapEvidence, &score)
			addErrorsIgnoreResults(results, author, ignoreEvidence)
		}
	}

	//results["error_wrap"][0] = errWrappScore.score / errWrappScore.fileLengh

	return results, nil
}

func getRelativePath(repoPath, filePath string) (string, error) {
	relativePath, err := filepath.Rel(repoPath, filePath)
	if err != nil {
		return "", err
	}
	return relativePath, nil
}

func addPanicUsageResults(results map[string][]Result, author string, panicEvidence []analysis.Evidence) {
	if len(panicEvidence) > 0 {
		addEvidence(results, "panic_usage", author, "2", panicEvidence)
	}
}

func addErrorWrapResults(results map[string][]Result,
	author string,
	wrapEvidence analysis.FileAnalysisResult,
	score *int) {
	if wrapEvidence.IsValid {
		if *score == 0 {
			*score = wrapEvidence.Score
		} else if *score > wrapEvidence.Score {
			*score = wrapEvidence.Score
		}

		addEvidence(results, "error_wrap", author, fmt.Sprintf("%d", *score), wrapEvidence.Evidences)
	}
}

func addErrorsIgnoreResults(results map[string][]Result, author string, ignoreEvidence []analysis.Evidence) {
	if len(ignoreEvidence) > 0 {
		addEvidence(results, "errors_ignore", author, "2", ignoreEvidence)
	}
}

func addEvidence(results map[string][]Result, metricID, author string, score string, evidence []analysis.Evidence) {
	found := false
	for i, result := range results[metricID] {
		if result.GitAuthor == author {
			results[metricID][i].Evidence = append(results[metricID][i].Evidence, evidence...)
			results[metricID][i].Score = score
			found = true
			break
		}
	}

	if !found {
		results[metricID] = append(results[metricID], Result{
			MetricID:  metricID,
			GitAuthor: author,
			Score:     score,
			Evidence:  evidence,
		})
	}
}
