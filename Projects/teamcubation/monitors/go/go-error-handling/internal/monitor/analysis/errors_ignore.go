package analysis

import (
	"regexp"
	"strings"
)

func AnalyzeErrorsIgnore(fileContent, filePath string) []Evidence {
	var evidence []Evidence
	lines := strings.Split(fileContent, "\n")

	re := regexp.MustCompile(`_,\s*err\s*:=`)

	for i, line := range lines {
		if re.MatchString(line) {
			if i+1 < len(lines) && !strings.Contains(lines[i+1], "if err != nil") {
				evidence = append(evidence, Evidence{
					File: filePath,
					Line: i + 1,
				})
			}
		}
	}
	return evidence
}
