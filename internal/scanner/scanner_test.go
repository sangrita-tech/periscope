package scanner_test

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sangrita-tech/periscope/internal/matcher"
	"github.com/sangrita-tech/periscope/internal/scanner"
)

type call struct {
	kind   string
	path   string
	depth  int
	isLast bool
}

func touch(t *testing.T, p string) {
	t.Helper()
	assert.NoError(t, os.WriteFile(p, []byte("x"), 0o644))
}

func Test_Scanner_Walk_VisitsDirsThenFiles_Sorted_DepthAndIsLast(t *testing.T) {
	root := t.TempDir()

	aDir := filepath.Join(root, "a_dir")
	bDir := filepath.Join(root, "b_dir")
	assert.NoError(t, os.MkdirAll(aDir, 0o755))
	assert.NoError(t, os.MkdirAll(bDir, 0o755))

	touch(t, filepath.Join(aDir, "a.txt"))
	touch(t, filepath.Join(root, "z.txt"))
	touch(t, filepath.Join(root, "m.txt"))

	var calls []call

	h := scanner.Handlers{
		OnDir: func(path string, d os.DirEntry, depth int, isLast bool) error {
			calls = append(calls, call{
				kind:   "dir",
				path:   filepath.Base(path),
				depth:  depth,
				isLast: isLast,
			})
			return nil
		},
		OnFile: func(path string, d os.DirEntry, depth int, isLast bool) error {
			calls = append(calls, call{
				kind:   "file",
				path:   filepath.Base(path),
				depth:  depth,
				isLast: isLast,
			})
			return nil
		},
	}

	s := scanner.New(root, nil, h)
	assert.NoError(t, s.Walk())

	want := []call{
		{kind: "dir", path: "a_dir", depth: 1, isLast: false},
		{kind: "file", path: "a.txt", depth: 2, isLast: true},
		{kind: "dir", path: "b_dir", depth: 1, isLast: false},
		{kind: "file", path: "m.txt", depth: 1, isLast: false},
		{kind: "file", path: "z.txt", depth: 1, isLast: true},
	}

	assert.Equal(t, want, calls)
}

func Test_Scanner_Walk_SkipsMatchedPaths_AndDoesNotDescend(t *testing.T) {
	root := t.TempDir()

	aDir := filepath.Join(root, "a_dir")
	bDir := filepath.Join(root, "b_dir")
	assert.NoError(t, os.MkdirAll(aDir, 0o755))
	assert.NoError(t, os.MkdirAll(bDir, 0o755))

	touch(t, filepath.Join(aDir, "a.txt"))
	touch(t, filepath.Join(bDir, "b.txt"))
	touch(t, filepath.Join(root, "root.txt"))

	m := matcher.New([]string{"a_dir*"})

	var calls []call
	h := scanner.Handlers{
		OnDir: func(path string, d os.DirEntry, depth int, isLast bool) error {
			calls = append(calls, call{kind: "dir", path: filepath.Base(path)})
			return nil
		},
		OnFile: func(path string, d os.DirEntry, depth int, isLast bool) error {
			calls = append(calls, call{kind: "file", path: filepath.Base(path)})
			return nil
		},
	}

	s := scanner.New(root, m, h)
	assert.NoError(t, s.Walk())

	for _, c := range calls {
		assert.NotEqual(t, "a_dir", c.path)
		assert.NotEqual(t, "a.txt", c.path)
	}

	var names []string
	for _, c := range calls {
		names = append(names, c.path)
	}
	assert.Contains(t, names, "b_dir")
	assert.Contains(t, names, "b.txt")
	assert.Contains(t, names, "root.txt")
}

func Test_Scanner_Walk_PropagatesHandlerError_StopsEarly(t *testing.T) {
	root := t.TempDir()

	assert.NoError(t, os.MkdirAll(filepath.Join(root, "a_dir"), 0o755))
	touch(t, filepath.Join(root, "a_dir", "a.txt"))
	touch(t, filepath.Join(root, "b.txt"))

	sentinel := errors.New("boom")

	var fileCalls int
	h := scanner.Handlers{
		OnFile: func(path string, d os.DirEntry, depth int, isLast bool) error {
			fileCalls++
			return sentinel
		},
	}

	s := scanner.New(root, nil, h)
	err := s.Walk()

	assert.ErrorIs(t, err, sentinel)
	assert.Equal(t, 1, fileCalls, "должны остановиться на первой ошибке")
}

func Test_Scanner_Walk_UsesSlashNormalizedRelativePathsForMatching(t *testing.T) {
	root := t.TempDir()

	nested := filepath.Join(root, "dir", "sub")
	assert.NoError(t, os.MkdirAll(nested, 0o755))
	touch(t, filepath.Join(nested, "f.txt"))

	m := matcher.New([]string{"dir/sub*"})

	var calls []call
	h := scanner.Handlers{
		OnDir: func(path string, d os.DirEntry, depth int, isLast bool) error {
			calls = append(calls, call{kind: "dir", path: filepath.Base(path)})
			return nil
		},
		OnFile: func(path string, d os.DirEntry, depth int, isLast bool) error {
			calls = append(calls, call{kind: "file", path: filepath.Base(path)})
			return nil
		},
	}

	s := scanner.New(root, m, h)
	assert.NoError(t, s.Walk())

	var names []string
	for _, c := range calls {
		names = append(names, c.path)
	}
	assert.Contains(t, names, "dir")
	assert.NotContains(t, names, "sub")
	assert.NotContains(t, names, "f.txt")

	_ = runtime.GOOS
}
