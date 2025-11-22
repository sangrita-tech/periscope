package scanner

import (
	"bufio"
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	contentbuilder "github.com/sangrita-tech/periscope/internal/content_builder"
	"github.com/sangrita-tech/periscope/internal/matcher"
	"github.com/sangrita-tech/periscope/internal/preprocess"
)

type Scanner struct {
	root           string
	pathMatcher    *matcher.Matcher
	contentMatcher *matcher.Matcher
	chain          *preprocess.Chain
	builder        *contentbuilder.ContentBuilder
}

func New(
	root string,
	pathMatcher *matcher.Matcher,
	contentMatcher *matcher.Matcher,
	chain *preprocess.Chain,
	builder *contentbuilder.ContentBuilder,
) *Scanner {
	return &Scanner{
		root:           root,
		pathMatcher:    pathMatcher,
		contentMatcher: contentMatcher,
		chain:          chain,
		builder:        builder,
	}
}

func (s *Scanner) Walk() error {
	return filepath.WalkDir(s.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path != s.root && s.pathMatcher != nil && !s.pathMatcher.Match(path) {
				return filepath.SkipDir
			}
			return nil
		}

		if s.pathMatcher != nil && !s.pathMatcher.Match(path) {
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

		if s.contentMatcher != nil {
			sc := bufio.NewScanner(strings.NewReader(content))
			for sc.Scan() {
				line := sc.Text()
				if s.contentMatcher.Match(line) {
					return nil
				}
			}
			if err := sc.Err(); err != nil {
				return err
			}
		}

		if s.chain != nil {
			processed, _, err := s.chain.Process(path, content)
			if err != nil {
				return err
			}
			content = processed
		}

		if strings.TrimSpace(content) == "" {
			return nil
		}

		if s.builder != nil {
			absPath, err := filepath.Abs(path)
			if err != nil {
				absPath = path
			}
			return s.builder.AddBlock("[FILE] "+absPath, content)
		}

		return nil
	})
}
