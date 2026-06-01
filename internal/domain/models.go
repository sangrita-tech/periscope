package domain

import "io/fs"

type Entry struct {
	Path    string
	RelPath string
	Data    []byte
	Meta    fs.FileInfo
}
