package walker

import (
	"io/fs"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/sangrita-tech/periscope/internal/domain"
)

type Walker struct{}

func New() *Walker {
	return &Walker{}
}

func (w *Walker) Walk(fsys fs.FS, root string) ([]domain.Entry, error) {
	root = path.Clean(root)

	meta, err := fs.Stat(fsys, root)
	if err != nil {
		return nil, err
	}

	if !meta.IsDir() {
		return w.walkSingleFile(fsys, root, meta)
	}

	entries, err := w.walkDirectory(fsys, root)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (w *Walker) walkSingleFile(
	fsys fs.FS,
	filePath string,
	meta fs.FileInfo,
) ([]domain.Entry, error) {
	data, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		return nil, err
	}

	return []domain.Entry{
		{
			Path:    filePath,
			RelPath: path.Clean(filePath),
			Data:    data,
			Meta:    meta,
		},
	}, nil
}

func (w *Walker) walkDirectory(fsys fs.FS, root string) ([]domain.Entry, error) {
	var entries []domain.Entry

	err := fs.WalkDir(fsys, root, func(currentPath string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if currentPath == root {
			return nil
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

		entries = append(entries, domain.Entry{
			Path:    currentPath,
			RelPath: path.Clean(currentPath),
			Data:    data,
			Meta:    meta,
		})

		log.Info().Str("path", currentPath).Msg("Found file")

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entries, nil
}
