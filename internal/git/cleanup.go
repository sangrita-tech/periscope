package git

import (
	"os"
	"path/filepath"
	"time"
)

func cleanupOldRepos(baseDir string, maxAge time.Duration) error {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}

	cutoff := time.Now().Add(-maxAge)

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		p := filepath.Join(baseDir, e.Name())
		info, err := e.Info()
		if err != nil {
			return err
		}

		if info.ModTime().Before(cutoff) {
			if err = os.RemoveAll(p); err != nil {
				return err
			}
		}
	}

	return nil
}
