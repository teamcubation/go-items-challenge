package analysis

import (
	"regexp"
	"strings"
)

func AnalyzePanicUsage(fileContent, filePath string) []Evidence {
	var evidence []Evidence
	if strings.Contains(filePath, "cmd") {
		return evidence
	}

	patterns := []string{
		`panic\(`,
		`log\.Fatal\(`,
		`log\.Fatalf\(`,
		`os\.Exit\(`,
	}

	combinedPattern := regexp.MustCompile(strings.Join(patterns, "|"))

	lines := strings.Split(fileContent, "\n")
	for i, line := range lines {
		if combinedPattern.MatchString(line) {
			evidence = append(evidence, Evidence{
				File: filePath,
				Line: i + 1,
			})
		}
	}

	return evidence
}
