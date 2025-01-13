package file

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func GetAllGoFiles(repoPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(repoPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".go") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("GetAllGoFiles: failed to read repo %s: %w", repoPath, err)
	}

	return files, nil
}
