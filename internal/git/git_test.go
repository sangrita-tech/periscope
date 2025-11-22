package git

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sangrita-tech/periscope/internal/config"
)

func Test_Fetch_EmptyRepoURL_ReturnsError(t *testing.T) {
	restore := withFakeGitSuccessOnly(t)
	defer restore()

	gc := New(&config.Git{CacheRoot: "cache_fetch_empty", MaxRepoAge: time.Hour})

	_, err := gc.Fetch("   ", "main")

	assert.Error(t, err)
}

func Test_Fetch_ClonesIntoCache_ReturnsCacheDir(t *testing.T) {
	restore := withFakeGitSuccessOnly(t)
	defer restore()

	cfg := &config.Git{CacheRoot: "cache_fetch_clone", MaxRepoAge: time.Hour}
	gc := New(cfg)

	repo := "https://example.com/a/b.git"
	branch := "main"

	out, err := gc.Fetch(repo, branch)

	assert.NoError(t, err)
	assert.DirExists(t, out)

	base := filepath.Join(os.TempDir(), cfg.CacheRoot)
	assert.Contains(t, out, base)
}

func Test_Fetch_WhenCacheExists_DoesNotCallGit_ReturnsExistingDir(t *testing.T) {
	restore := withFakeGitFailOnly(t)
	defer restore()

	cfg := &config.Git{CacheRoot: "cache_fetch_exists", MaxRepoAge: time.Hour}
	gc := New(cfg)

	repo := "https://example.com/a/b.git"
	branch := "main"

	base := filepath.Join(os.TempDir(), cfg.CacheRoot)
	cacheDir := buildCacheDir(base, repo, branch)
	assert.NoError(t, os.MkdirAll(cacheDir, 0o755))

	out, err := gc.Fetch(repo, branch)

	assert.NoError(t, err)
	assert.Equal(t, cacheDir, out)
}

func withFakeGitSuccessOnly(t *testing.T) func() {
	t.Helper()

	dir := t.TempDir()

	name := "git"
	if runtime.GOOS == "windows" {
		name = "git.bat"
	}

	path := filepath.Join(dir, name)
	script := fakeGitScriptSuccess(runtime.GOOS)

	assert.NoError(t, os.WriteFile(path, []byte(script), 0o755))

	// ВАЖНО: заменяем PATH полностью, а не добавляем!
	oldPath := os.Getenv("PATH")
	assert.NoError(t, os.Setenv("PATH", dir))

	return func() { _ = os.Setenv("PATH", oldPath) }
}

func withFakeGitFailOnly(t *testing.T) func() {
	t.Helper()

	dir := t.TempDir()

	name := "git"
	if runtime.GOOS == "windows" {
		name = "git.bat"
	}

	path := filepath.Join(dir, name)
	script := fakeGitScriptFail(runtime.GOOS)

	assert.NoError(t, os.WriteFile(path, []byte(script), 0o755))

	oldPath := os.Getenv("PATH")
	assert.NoError(t, os.Setenv("PATH", dir))

	return func() { _ = os.Setenv("PATH", oldPath) }
}

func fakeGitScriptSuccess(goos string) string {
	if goos == "windows" {
		return "@echo off\r\n" +
			"if \"%1\"==\"--version\" exit /b 0\r\n" +
			"set target=\r\n" +
			"for %%x in (%*) do set target=%%x\r\n" +
			"mkdir \"%target%\" >nul 2>nul\r\n" +
			"exit /b 0\r\n"
	}

	return "#!/bin/sh\n" +
		"if [ \"$1\" = \"--version\" ]; then exit 0; fi\n" +
		"eval \"target=\\${$#}\"\n" +
		"mkdir -p \"$target\"\n" +
		"exit 0\n"
}

func fakeGitScriptFail(goos string) string {
	if goos == "windows" {
		return "@echo off\r\n" +
			"if \"%1\"==\"--version\" exit /b 0\r\n" +
			"exit /b 1\r\n"
	}

	return "#!/bin/sh\n" +
		"if [ \"$1\" = \"--version\" ]; then exit 0; fi\n" +
		"exit 1\n"
}
