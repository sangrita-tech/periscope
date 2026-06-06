package walker

import (
	"io/fs"
	"path"
	"strings"

	"github.com/sangrita-tech/periscope/internal/ignore"
	"github.com/sangrita-tech/periscope/internal/model"
)

type Walker struct {
	ignore *ignore.Matcher
}

func New(matcher *ignore.Matcher) *Walker {
	return &Walker{
		ignore: matcher,
	}
}

func (w *Walker) Walk(src model.Source) ([]model.Entry, error) {
	var entries []model.Entry

	root := path.Clean(src.Root)

	err := fs.WalkDir(src.Fsys, root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		currentPath := path.Clean(p)
		relPath := relativePath(root, currentPath)

		if w.ignore.Match(relPath) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		meta, err := d.Info()
		if err != nil {
			return err
		}

		data, err := fs.ReadFile(src.Fsys, currentPath)
		if err != nil {
			return err
		}

		entries = append(entries, model.Entry{
			Path:    currentPath,
			RelPath: relPath,
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

func relativePath(root, currentPath string) string {
	if currentPath == root {
		return "."
	}

	return strings.TrimPrefix(currentPath, root+"/")
}
