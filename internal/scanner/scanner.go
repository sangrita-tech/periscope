package scanner

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/sangrita-tech/periscope/internal/matcher"
	"github.com/sangrita-tech/periscope/internal/transformer"
)

func WalkProcessedFiles(root string, matcher *matcher.Matcher, pipeline *transformer.Pipeline, handler func(path, content string) error) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if matcher.ShouldIgnorePath(path) {
				return filepath.SkipDir
			}
			return nil
		}

		if matcher.ShouldIgnorePath(path) {
			return nil
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if bytes.IndexByte(contentBytes, 0x00) != -1 {
			return nil
		}

		content := string(contentBytes)

		if matcher.ShouldIgnoreContent(content) {
			return nil
		}

		processed, _, err := pipeline.Process(path, content)
		if err != nil {
			return err
		}

		if len(bytes.TrimSpace([]byte(processed))) == 0 {
			return nil
		}

		return handler(path, processed)
	})
}
