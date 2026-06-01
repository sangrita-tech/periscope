package walker

import (
	"io/fs"
	"path"

	"github.com/sangrita-tech/periscope/internal/domain"
)

type Walker struct{}

func New() *Walker {
	return &Walker{}
}

func (w *Walker) Walk(fsys fs.FS, root string) ([]domain.Entry, error) {
	root = path.Clean(root)

	var entries []domain.Entry

	err := fs.WalkDir(fsys, root, func(currentPath string, dirEntry fs.DirEntry, walkErr error) error {
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

		data, err := fs.ReadFile(fsys, currentPath)
		if err != nil {
			return err
		}

		currentPath = path.Clean(currentPath)

		entries = append(entries, domain.Entry{
			Path:    currentPath,
			RelPath: currentPath,
			Data:    data,
			Meta:    meta,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entries, nil
}
