package source

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
	"github.com/sangrita-tech/periscope/internal/domain"
)

const root = "."

func ResolveSource(target string) (domain.Source, error) {
	target = strings.TrimSpace(target)
	if target == "" {
		target = "."
	}

	if isSourceRemote(target) {
		return resolveGitSource(target)
	}

	return resolveLocalSource(target)
}

func isSourceRemote(target string) bool {
	return strings.HasPrefix(target, "git+") ||
		strings.HasPrefix(target, "https://") ||
		strings.HasPrefix(target, "http://") ||
		strings.HasPrefix(target, "ssh://") ||
		(strings.Contains(target, "@") && strings.Contains(target, ":"))
}

func resolveLocalSource(target string) (domain.Source, error) {
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return domain.Source{}, fmt.Errorf("absolute path: %w", err)
	}

	if _, err := os.Stat(absTarget); err != nil {
		return domain.Source{}, fmt.Errorf("stat local path: %w", err)
	}

	parent := filepath.Dir(absTarget)
	base := filepath.Base(absTarget)

	return domain.Source{
		Fsys: os.DirFS(parent),
		Root: filepath.ToSlash(base),
	}, nil
}

func resolveGitSource(target string) (domain.Source, error) {
	mux := fsimpl.NewMux()
	mux.Add(gitfs.FS)

	fsys, err := mux.Lookup("git+https://github.com/hairyhenderson/go-fsimpl")
	if err != nil {
		return domain.Source{}, fmt.Errorf("lookup git filesystem: %w", err)
	}

	return domain.Source{
		Fsys: fsys,
		Root: root,
	}, nil
}
