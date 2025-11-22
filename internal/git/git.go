package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sangrita-tech/periscope/internal/config"
)

const gitBinary = "git"

type Git struct {
	cfg *config.Git
}

func New(cfg *config.Git) *Git {
	return &Git{cfg: cfg}
}

func (g *Git) Fetch(repoURL, branch string) (string, error) {
	if !isGitInstalled() {
		return "", fmt.Errorf("git is not installed")
	}

	repoURL = strings.TrimSpace(repoURL)
	if repoURL == "" {
		return "", fmt.Errorf("empty repository url")
	}

	cacheBase := cacheRootDir(g.cfg.CacheRoot)
	cacheDir := buildCacheDir(cacheBase, repoURL, branch)

	if err := os.MkdirAll(filepath.Dir(cacheDir), 0o755); err != nil {
		return "", fmt.Errorf("create cache dir: %w", err)
	}

	if info, err := os.Stat(cacheDir); err == nil && info.IsDir() {
		return cacheDir, nil
	}

	args := []string{"clone", "--depth=1"}
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	args = append(args, repoURL, cacheDir)

	cmd := exec.Command(gitBinary, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if err = os.RemoveAll(cacheDir); err != nil {
			return "", err
		}
		return "", fmt.Errorf("git clone failed: %w: %s", err, string(out))
	}

	if err = cleanupOldRepos(cacheBase, g.cfg.MaxRepoAge); err != nil {
		return "", err
	}
	return cacheDir, nil
}

func isGitInstalled() bool {
	return exec.Command(gitBinary, "--version").Run() == nil
}
