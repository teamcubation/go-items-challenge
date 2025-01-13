package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

func GetLastCommitHash(filePath, repoPath string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("error opening git repository: %w", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("work tree is not tracked by Git in the repository: %w", err)
	}

	_, err = wt.Filesystem.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("the file '%s' is not tracked by Git in the repository", filePath)
	}

	commits, err := repo.Log(&git.LogOptions{FileName: &filePath})
	if err != nil {
		return "", fmt.Errorf("error getting commits: %w", err)
	}

	if commits != nil {
		commit, err := commits.Next()
		if err != nil {
			return "", fmt.Errorf("error getting the latest commit: %w", err)
		}
		return commit.Hash.String(), nil
	}

	return "", nil
}
