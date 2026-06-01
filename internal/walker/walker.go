package walker

import (
	"io/fs"
	"path"
	"strings"

	"github.com/sangrita-tech/periscope/internal/domain"
	"github.com/sangrita-tech/periscope/internal/ignore"
)

type Walker struct {
	ignore *ignore.Matcher
}

func New(matcher *ignore.Matcher) *Walker {
	return &Walker{
		ignore: matcher,
	}
}

func (w *Walker) Walk(source domain.Source) ([]domain.Entry, error) {
	var entries []domain.Entry

	err := fs.WalkDir(source.Fsys, source.Root, func(currentPath string, dirEntry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		currentPath = path.Clean(currentPath)
		relPath := makeRelPath(source.Root, currentPath)

		ignored := w.shouldIgnore(relPath, dirEntry.IsDir())
		if ignored {
			if dirEntry.IsDir() {
				return fs.SkipDir
			}

			return nil
		}

		if dirEntry.IsDir() {
			return nil
		}

		meta, err := dirEntry.Info()
		if err != nil {
			return err
		}

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

func (w *Walker) shouldIgnore(relPath string, isDir bool) bool {
	return w.ignore != nil && w.ignore.Match(relPath, isDir)
}

func makeRelPath(root, currentPath string) string {
	root = cleanWalkPath(root)
	currentPath = cleanWalkPath(currentPath)

	if root == "" || root == "." {
		return currentPath
	}

	if currentPath == root {
		return ""
	}

	if relPath, ok := strings.CutPrefix(currentPath, root+"/"); ok {
		return relPath
	}

	return strings.TrimLeft(strings.TrimPrefix(currentPath, root), "/")
}

func cleanWalkPath(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "\\", "/")
	value = strings.TrimLeft(value, "/")
	value = path.Clean(value)

	if value == "." {
		return ""
	}

	return value
}
