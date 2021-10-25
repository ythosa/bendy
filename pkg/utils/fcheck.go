package utils

import (
	"fmt"
	"os"
)

func CheckIsFilePathsValid(paths []string) error {
	for _, fp := range paths {
		if _, err := os.Stat(fp); err != nil {
			return fmt.Errorf("error while getting file stat: %w", err)
		}
	}

	return nil
}
