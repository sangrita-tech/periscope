package matcher

import (
	"bufio"
	"path/filepath"
	"strings"
)

type Matcher struct {
	ignorePaths    []string
	ignoreContents []string
}

func New(ignorePaths, ignoreContents []string) *Matcher {
	return &Matcher{
		ignorePaths:    ignorePaths,
		ignoreContents: ignoreContents,
	}
}

func (m *Matcher) ShouldIgnorePath(path string) bool {
	if len(m.ignorePaths) == 0 {
		return false
	}

	normPath := filepath.ToSlash(path)

	for _, p := range m.ignorePaths {
		if p == "" {
			continue
		}

		p = filepath.ToSlash(p)

		if ok, _ := filepath.Match(p, normPath); ok {
			return true
		}

		if ok, _ := filepath.Match(p, filepath.Base(normPath)); ok {
			return true
		}

		if strings.Contains(normPath, p) {
			return true
		}
	}

	return false
}

func (m *Matcher) ShouldIgnoreContent(content string) bool {
	if len(m.ignoreContents) == 0 {
		return false
	}

	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()

		for _, pattern := range m.ignoreContents {
			pattern = strings.TrimSpace(pattern)
			if pattern == "" {
				continue
			}

			if !strings.ContainsAny(pattern, "*?[]") {
				if strings.Contains(line, pattern) {
					return true
				}
				continue
			}

			if ok, _ := filepath.Match(pattern, line); ok {
				return true
			}
		}
	}

	return false
}
