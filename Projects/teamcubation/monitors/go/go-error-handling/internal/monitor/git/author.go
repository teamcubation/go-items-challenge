package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetGitAuthor(commitHash string, repoPath string) (string, error) {
	cmd := exec.Command("git", "show", "-s", "--format=%an <%ae>", commitHash)
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "Unknown", fmt.Errorf("error getting git author: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
