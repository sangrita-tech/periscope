package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/sangrita-tech/periscope/internal/matcher"
)

type Handlers struct {
	OnDir  func(path string, d fs.DirEntry, depth int, isLast bool) error
	OnFile func(path string, d fs.DirEntry, depth int, isLast bool) error
}

type Scanner struct {
	root        string
	pathMatcher *matcher.Matcher
	handlers    Handlers
}

func New(root string, pathMatcher *matcher.Matcher, handlers Handlers) *Scanner {
	return &Scanner{
		root:        root,
		pathMatcher: pathMatcher,
		handlers:    handlers,
	}
}

func (s *Scanner) Walk() error {
	absRoot, err := filepath.Abs(s.root)
	if err != nil {
		absRoot = s.root
	}
	return s.walkDir(absRoot, 0)
}

func (s *Scanner) walkDir(dir string, depth int) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	dirs := make([]fs.DirEntry, 0, len(entries))
	files := make([]fs.DirEntry, 0, len(entries))

	for _, e := range entries {
		full := filepath.Join(dir, e.Name())

		if s.pathMatcher != nil && !s.pathMatcher.Match(full) {
			continue
		}

		if e.IsDir() {
			dirs = append(dirs, e)
		} else {
			files = append(files, e)
		}
	}

	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	all := append(dirs, files...)

	for i, e := range all {
		full := filepath.Join(dir, e.Name())
		isLast := i == len(all)-1

		if e.IsDir() {
			if s.handlers.OnDir != nil {
				if err := s.handlers.OnDir(full, e, depth+1, isLast); err != nil {
					return err
				}
			}
			if err := s.walkDir(full, depth+1); err != nil {
				return err
			}
			continue
		}

		if s.handlers.OnFile != nil {
			if err := s.handlers.OnFile(full, e, depth+1, isLast); err != nil {
				return err
			}
		}
	}

	return nil
}
