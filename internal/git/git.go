package git

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sangrita-tech/periscope/internal/config"
)

const gitBinary = "git"

type Git struct {
	cfg *config.Git
}

func New(cfg *config.Git) *Git {
	return &Git{cfg: cfg}
}

func (g *Git) CloneRepo(repo, branch string) (string, error) {
	if !isGitInstalled() {
		return "", errors.New("git is not installed")
	}

	normalizedRepo := normalizeRepoURL(repo)

	cacheBase := filepath.Join(os.TempDir(), g.cfg.CacheRoot)
	cacheDir := buildCachePath(cacheBase, normalizedRepo, branch)

	parent := filepath.Dir(cacheDir)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return "", errors.New("failed to create cache directory")
	}

	if info, err := os.Stat(cacheDir); err == nil && info.IsDir() {
		return cacheDir, nil
	}

	args := []string{"clone", "--depth=1"}
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	args = append(args, normalizedRepo, cacheDir)

	cmd := exec.Command(gitBinary, args...)
	if _, err := cmd.CombinedOutput(); err != nil {
		_ = os.RemoveAll(cacheDir)
		return "", errors.New("failed to clone repository")
	}

	_ = cleanupOldRepos(cacheBase, g.cfg.MaxRepoAge)

	return cacheDir, nil
}

func isGitInstalled() bool {
	cmd := exec.Command(gitBinary, "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func buildCachePath(cacheRoot, repoURL, branch string) string {
	provider, owner, project := parseRepoParts(repoURL)
	if provider == "" {
		provider = "unknown"
	}
	if owner == "" {
		owner = "unknown"
	}
	if project == "" {
		project = "unknown"
	}

	dirName := fmt.Sprintf("%s_%s_%s", provider, owner, project)
	if branch != "" {
		dirName = fmt.Sprintf("%s_%s", dirName, sanitizeBranch(branch))
	}

	return filepath.Join(cacheRoot, dirName)
}

func parseRepoParts(repoURL string) (provider, owner, project string) {
	if strings.HasPrefix(repoURL, "git@") {
		parts := strings.SplitN(repoURL, ":", 2)
		if len(parts) == 2 {
			hostPart := parts[0]
			pathPart := parts[1]

			if strings.Contains(hostPart, "github.com") {
				provider = "github"
			} else if strings.Contains(hostPart, "gitlab.com") {
				provider = "gitlab"
			}

			segs := strings.Split(strings.TrimSuffix(pathPart, ".git"), "/")
			if len(segs) >= 2 {
				owner = segs[0]
				project = segs[1]
			}
			return
		}
	}

	u, err := url.Parse(repoURL)
	if err != nil || u.Host == "" {
		return "", "", ""
	}

	host := strings.ToLower(u.Host)
	if strings.Contains(host, "github.com") {
		provider = "github"
	} else if strings.Contains(host, "gitlab.com") {
		provider = "gitlab"
	}

	segs := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(segs) >= 2 {
		owner = segs[0]
		project = strings.TrimSuffix(segs[1], ".git")
	}

	return
}

func sanitizeBranch(branch string) string {
	b := branch
	replacements := []string{"/", "\\", ":", "@", "?", "&", "=", ".", " "}
	for _, r := range replacements {
		b = strings.ReplaceAll(b, r, "_")
	}
	return b
}

func normalizeRepoURL(raw string) string {
	if strings.HasPrefix(raw, "git@") {
		return raw
	}

	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return raw
	}

	host := strings.ToLower(u.Host)
	if !strings.Contains(host, "github.com") && !strings.Contains(host, "gitlab.com") {
		return raw
	}

	path := strings.TrimSuffix(u.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return raw
	}

	owner := parts[1]
	repo := parts[2]

	if !strings.HasSuffix(repo, ".git") {
		repo = repo + ".git"
	}

	u.Path = fmt.Sprintf("/%s/%s", owner, repo)
	u.RawQuery = ""
	u.Fragment = ""

	return u.String()
}

func cleanupOldRepos(baseDir string, maxAge time.Duration) error {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return errors.New("failed to read cache directory")
	}

	cutoff := time.Now().Add(-maxAge)

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		p := filepath.Join(baseDir, e.Name())
		info, err := e.Info()
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return errors.New("failed to inspect repository directory")
		}

		if info.ModTime().Before(cutoff) {
			if err = os.RemoveAll(p); err != nil {
				return errors.New("failed to remove old repository")
			}
		}
	}

	return nil
}
