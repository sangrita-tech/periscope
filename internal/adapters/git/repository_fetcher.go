package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sangrita-tech/periscope/internal/config"
	"github.com/sangrita-tech/periscope/internal/domain"
)

const binaryName = "git"
const defaultCacheRoot = "~/.periscope/git_cache"
const defaultMaxRepoAge = 14 * 24 * time.Hour
const cacheDirMode = 0o755

type repositoryFetcher struct {
	cfg *config.Git
}

func NewRepositoryFetcher(cfg *config.Git) *repositoryFetcher {
	if cfg.CacheRoot == "" {
		cfg.CacheRoot = defaultCacheRoot
	}
	if cfg.MaxRepoAge == 0 {
		cfg.MaxRepoAge = defaultMaxRepoAge
	}
	return &repositoryFetcher{
		cfg: cfg,
	}
}

func (f *repositoryFetcher) Fetch(repository domain.Repository) (string, error) {
	if err := cleanupOldRepos(); err != nil {
		return "", fmt.Errorf("failed to cleanup old repositories: %w", err)
	}

	if repository.URL == "" {
		return "", errors.New("repository url is required")
	}

	if !isGitInstalled() {
		return "", errors.New("git is not installed")
	}

	repoCacheDirPath := buildRepoCacheDirPath(repository)
	if isDirExists(repoCacheDirPath) {
		return repoCacheDirPath, nil
	}

	if err := os.MkdirAll(filepath.Dir(repoCacheDirPath), cacheDirMode); err != nil {
		return "", err
	}

	err := cloneRepo(repository, repoCacheDirPath)
	if err != nil {
		err = os.RemoveAll(repoCacheDirPath)
		if err != nil {
			return "", fmt.Errorf("failed to clone repository and cleanup cache dir: %w", err)
		}
		return "", errors.New("failed to clone repository")
	}

	return repoCacheDirPath, nil
}

func cleanupOldRepos() error {
	entries, err := os.ReadDir(defaultCacheRoot)
	if err != nil {
		return err
	}

	cutoff := time.Now().Add(-defaultMaxRepoAge)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		if info.ModTime().Before(cutoff) {
			if err = os.RemoveAll(filepath.Join(defaultCacheRoot, entry.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}

func isGitInstalled() bool {
	return exec.Command(binaryName, "--version").Run() == nil
}

func buildRepoCacheDirPath(repository domain.Repository) string {
	dir := sanitizeString(string(repository.Provider)) + "_" + sanitizeString(repository.Owner) + "_" + sanitizeString(repository.Name)
	if repository.Branch != "" {
		dir += "_" + sanitizeString(repository.Branch)
	}
	return filepath.Join(defaultCacheRoot, dir)
}

func sanitizeString(value string) string {
	value = strings.TrimSpace(value)
	replacements := []string{
		"/", "\\", ":", "@", "?", "&", "=", ".", " ", "%", "#",
		"*", "\"", "<", ">", "|", "\t", "\n", "\r",
	}

	for _, replacement := range replacements {
		value = strings.ReplaceAll(value, replacement, "_")
	}

	for strings.Contains(value, "__") {
		value = strings.ReplaceAll(value, "__", "_")
	}

	return strings.Trim(value, "_")
}

func isDirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func cloneRepo(repository domain.Repository, repoCacheDirPath string) error {
	args := []string{"clone", "--depth=1"}
	if repository.Branch != "" {
		args = append(args, "--branch", repository.Branch)
	}
	args = append(args, repository.URL, repoCacheDirPath)

	command := exec.Command(binaryName, args...)

	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %w: %s", err, string(output))
	}

	return nil
}
