package git

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func cacheRootDir(cacheRoot string) string {
	return filepath.Join(os.TempDir(), cacheRoot)
}

func buildCacheDir(cacheBase, repoURL, branch string) string {
	provider, owner, repo := parseRepoParts(repoURL)

	if provider == "" {
		provider = "unknown"
	}
	if owner == "" {
		owner = "unknown"
	}
	if repo == "" {
		repo = "unknown"
	}

	provider = sanitizePart(provider)
	owner = sanitizePart(owner)
	repo = sanitizePart(repo)

	dir := provider + "_" + owner + "_" + repo
	if branch != "" {
		dir += "_" + sanitizeBranch(branch)
	}

	return filepath.Join(cacheBase, dir)
}

func parseRepoParts(repoURL string) (provider, owner, repo string) {
	repoURL = strings.TrimSpace(repoURL)
	if repoURL == "" {
		return "", "", ""
	}

	if strings.HasPrefix(repoURL, "git@") {
		parts := strings.SplitN(repoURL, ":", 2)
		if len(parts) != 2 {
			return "", "", ""
		}

		hostPart := strings.ToLower(parts[0])
		pathPart := strings.TrimSuffix(parts[1], ".git")
		segs := strings.Split(strings.Trim(pathPart, "/"), "/")

		if strings.Contains(hostPart, "github.com") {
			provider = "github"
		} else if strings.Contains(hostPart, "gitlab.com") {
			provider = "gitlab"
		}

		if len(segs) >= 2 {
			repo = segs[len(segs)-1]
			owner = strings.Join(segs[:len(segs)-1], "_")
		}

		return
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
		repo = strings.TrimSuffix(segs[len(segs)-1], ".git")
		owner = strings.Join(segs[:len(segs)-1], "_")
	}

	return
}

func sanitizeBranch(branch string) string {
	return sanitizePart(branch)
}

func sanitizePart(s string) string {
	s = strings.TrimSpace(s)
	replacements := []string{
		"/", "\\", ":", "@", "?", "&", "=", ".", " ", "%", "#",
		"*", "\"", "<", ">", "|", "\t", "\n", "\r",
	}

	for _, r := range replacements {
		s = strings.ReplaceAll(s, r, "_")
	}

	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}

	s = strings.Trim(s, "_")

	if s == "" {
		return "unknown"
	}

	return s
}
