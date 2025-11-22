package git

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_cleanupOldRepos_RemovesOldDirectories(t *testing.T) {
	base := t.TempDir()

	oldDir := filepath.Join(base, "old")
	newDir := filepath.Join(base, "new")

	assert.NoError(t, os.MkdirAll(oldDir, 0o755))
	assert.NoError(t, os.MkdirAll(newDir, 0o755))

	oldTime := time.Now().Add(-48 * time.Hour)
	newTime := time.Now().Add(-1 * time.Hour)

	assert.NoError(t, os.Chtimes(oldDir, oldTime, oldTime))
	assert.NoError(t, os.Chtimes(newDir, newTime, newTime))

	err := cleanupOldRepos(base, 24*time.Hour)

	assert.NoError(t, err)
	_, errOld := os.Stat(oldDir)
	_, errNew := os.Stat(newDir)

	assert.True(t, os.IsNotExist(errOld))
	assert.NoError(t, errNew)
}

func Test_cleanupOldRepos_IgnoresFiles(t *testing.T) {
	base := t.TempDir()

	f := filepath.Join(base, "file.txt")
	assert.NoError(t, os.WriteFile(f, []byte("x"), 0o644))

	err := cleanupOldRepos(base, time.Hour)

	assert.NoError(t, err)
	_, statErr := os.Stat(f)
	assert.NoError(t, statErr)
}

func Test_cleanupOldRepos_NonExistingBase_ReturnsError(t *testing.T) {
	err := cleanupOldRepos(filepath.Join(t.TempDir(), "nope"), time.Hour)
	assert.Error(t, err)
}
