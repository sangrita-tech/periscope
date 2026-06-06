package source

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
	"github.com/sangrita-tech/periscope/internal/model"
)

const root = "."

func ResolveSource(target string) (model.Source, error) {
	target = strings.TrimSpace(target)
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

func resolveLocalSource(target string) (model.Source, error) {
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return model.Source{}, fmt.Errorf("absolute path: %w", err)
	}

	if _, err := os.Stat(absTarget); err != nil {
		return model.Source{}, fmt.Errorf("stat local path: %w", err)
	}

	parent := filepath.Dir(absTarget)
	base := filepath.Base(absTarget)

	return model.Source{
		Fsys: os.DirFS(parent),
		Root: filepath.ToSlash(base),
		Name: base,
	}, nil
}

func resolveGitSource(target string) (model.Source, error) {
	mux := fsimpl.NewMux()
	mux.Add(gitfs.FS)

	gitTarget := target
	if !strings.HasPrefix(gitTarget, "git+") {
		gitTarget = "git+" + gitTarget
	}

	fsys, err := mux.Lookup(gitTarget)
	if err != nil {
		return model.Source{}, fmt.Errorf("lookup git filesystem: %w", err)
	}

	sourceName := strings.TrimSpace(target)
	sourceName = strings.TrimPrefix(sourceName, "git+")
	sourceName = strings.TrimSuffix(sourceName, ".git")
	sourceName = path.Base(sourceName)

	return model.Source{
		Fsys: fsys,
		Root: root,
		Name: sourceName,
	}, nil
}
