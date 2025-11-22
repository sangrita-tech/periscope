package scanner_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	contentbuilder "github.com/sangrita-tech/periscope/internal/content_builder"
	"github.com/sangrita-tech/periscope/internal/matcher"
	"github.com/sangrita-tech/periscope/internal/preprocess"
	"github.com/sangrita-tech/periscope/internal/scanner"
)

func Test_Walk_NoMatchersNoChain_AddsAllTextFiles(t *testing.T) {
	root := t.TempDir()

	assert.NoError(t, os.WriteFile(filepath.Join(root, "a.txt"), []byte("hello"), 0o644))
	assert.NoError(t, os.WriteFile(filepath.Join(root, "b.txt"), []byte("world\n"), 0o644))

	b := contentbuilder.New()
	s := scanner.New(root, nil, nil, nil, b)

	err := s.Walk()

	assert.NoError(t, err)
	out := b.Result()
	assert.Contains(t, out, "hello")
	assert.Contains(t, out, "world")
	assert.Contains(t, out, "[LABEL]")
}

func Test_Walk_PathMatcher_SkipsNonMatchingFilesAndDirs(t *testing.T) {
	root := t.TempDir()

	keepDir := filepath.Join(root, "keep")
	skipDir := filepath.Join(root, "skip")

	assert.NoError(t, os.MkdirAll(keepDir, 0o755))
	assert.NoError(t, os.MkdirAll(skipDir, 0o755))

	assert.NoError(t, os.WriteFile(filepath.Join(keepDir, "ok.txt"), []byte("ok"), 0o644))
	assert.NoError(t, os.WriteFile(filepath.Join(skipDir, "no.txt"), []byte("no"), 0o644))

	pm := matcher.New([]string{"keep"})
	b := contentbuilder.New()
	s := scanner.New(root, pm, nil, nil, b)

	err := s.Walk()

	assert.NoError(t, err)
	out := b.Result()
	assert.Contains(t, out, "ok")
	assert.NotContains(t, out, "no")
}

func Test_Walk_ContentMatcher_IgnoresFileIfAnyLineMatches(t *testing.T) {
	root := t.TempDir()

	assert.NoError(t, os.WriteFile(filepath.Join(root, "a.txt"), []byte("hello\nSECRET\nworld\n"), 0o644))
	assert.NoError(t, os.WriteFile(filepath.Join(root, "b.txt"), []byte("hello\nworld\n"), 0o644))

	cm := matcher.New([]string{"*SECRET*"})
	b := contentbuilder.New()
	s := scanner.New(root, nil, cm, nil, b)

	err := s.Walk()

	assert.NoError(t, err)
	out := b.Result()
	assert.NotContains(t, out, "SECRET")
	assert.Contains(t, out, "hello\nworld")
}

func Test_Walk_SkipsBinaryFiles(t *testing.T) {
	root := t.TempDir()

	bin := []byte{0x01, 0x02, 0x00, 0x03}
	assert.NoError(t, os.WriteFile(filepath.Join(root, "bin.dat"), bin, 0o644))
	assert.NoError(t, os.WriteFile(filepath.Join(root, "text.txt"), []byte("text"), 0o644))

	b := contentbuilder.New()
	s := scanner.New(root, nil, nil, nil, b)

	err := s.Walk()

	assert.NoError(t, err)
	out := b.Result()
	assert.NotContains(t, out, "bin.dat")
	assert.Contains(t, out, "text")
}

func Test_Walk_ChainProcessesContentBeforeAdding(t *testing.T) {
	root := t.TempDir()

	in := "// c1\n\n\nx\n\n"
	assert.NoError(t, os.WriteFile(filepath.Join(root, "a.txt"), []byte(in), 0o644))

	ch := preprocess.New().
		AddStripComments().
		AddCollapseEmptyLines()

	b := contentbuilder.New()
	s := scanner.New(root, nil, nil, ch, b)

	err := s.Walk()

	assert.NoError(t, err)

	out := b.Result()
	assert.Contains(t, out, "x\n")
	assert.NotContains(t, out, "// c1")

	parts := strings.SplitN(out, "\n\n", 2)
	assert.Len(t, parts, 2)
	body := parts[1]

	assert.NotContains(t, body, "\n\n\n")
}
