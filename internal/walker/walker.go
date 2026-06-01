package walker

import (
	"io/fs"
	"path"
	"strings"

	"github.com/sangrita-tech/periscope/internal/domain"
)

type Walker struct{}

func New() *Walker {
	return &Walker{}
}

func (w *Walker) Walk(source domain.Source) ([]domain.Entry, error) {
	var entries []domain.Entry

	err := fs.WalkDir(source.Fsys, source.Root, func(currentPath string, dirEntry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if dirEntry.IsDir() {
			return nil
		}

		meta, err := dirEntry.Info()
		if err != nil {
			return err
		}

		currentPath = path.Clean(currentPath)
		relPath := strings.TrimPrefix(path.Clean(strings.TrimPrefix(currentPath, source.Root)), "/")

		entry := domain.Entry{
			Path:    currentPath,
			RelPath: relPath,
			Meta:    meta,
		}

		data, err := fs.ReadFile(source.Fsys, currentPath)
		if err == nil {
			entry.Data = data
		}

		entries = append(entries, entry)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entries, nil
}
