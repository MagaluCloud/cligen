package common

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func GetFileDst(dst string, objectKey string) (string, error) {
	fileName := path.Base(objectKey)

	if dst == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %w", err)
		}

		return filepath.Join(cwd, fileName), nil
	}

	info, err := os.Stat(dst)
	if (err == nil && info.IsDir()) || dst[len(dst)-1] == os.PathSeparator {
		dst = filepath.Join(dst, fileName)
	}

	return dst, nil
}
