package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_cacheRootDir_AppendsTempDir_ReturnsJoinedPath(t *testing.T) {
	root := "periscope_cache_x"

	got := filepath.Clean(cacheRootDir(root))
	want := filepath.Clean(filepath.Join(os.TempDir(), root))

	assert.Equal(t, want, got)
}

func Test_buildCacheDir_SameInputs_ReturnSamePath(t *testing.T) {
	base := "/tmp/base"
	repo := "https://example.com/a/b.git"
	branch := "main"

	p1 := buildCacheDir(base, repo, branch)
	p2 := buildCacheDir(base, repo, branch)

	assert.Equal(t, p1, p2)
}

func Test_buildCacheDir_DifferentBranch_ReturnDifferentPath(t *testing.T) {
	base := "/tmp/base"
	repo := "https://example.com/a/b.git"

	p1 := buildCacheDir(base, repo, "main")
	p2 := buildCacheDir(base, repo, "dev")

	assert.NotEqual(t, p1, p2)
}

func Test_buildCacheDir_EmptyBranch_DoesNotAppendSanitizedBranch(t *testing.T) {
	base := "/tmp/base"
	repo := "https://example.com/a/b.git"

	p := buildCacheDir(base, repo, "")

	assert.Equal(t, filepath.Clean(base), filepath.Clean(filepath.Dir(p)))
	_, dir := filepath.Split(p)
	assert.NotContains(t, dir, "__")
	assert.NotContains(t, dir, "main")
}

func Test_sanitizeBranch_ReplacesUnsafeChars_ReturnsSafeString(t *testing.T) {
	got := sanitizeBranch(" feature/x:y@z?.go ")
	assert.Equal(t, "feature_x_y_z__go", got)
}

func Test_sanitizeBranch_Empty_ReturnsEmpty(t *testing.T) {
	assert.Equal(t, "", sanitizeBranch("   "))
}
