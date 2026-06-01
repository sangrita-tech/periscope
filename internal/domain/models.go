package domain

import "io/fs"

type Entry struct {
	Path    string
	RelPath string
	Data    []byte
	Meta    fs.FileInfo
}

type Source struct {
	Fsys fs.FS
	Root string
}
