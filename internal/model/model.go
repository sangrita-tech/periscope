package model

import "io/fs"

type Entry struct {
	Path    string
	RelPath string
	Data    []byte
	Meta    fs.FileInfo
}

type Replacement struct {
	Pattern string `yaml:"pattern"`
	Value   string `yaml:"value"`
}

type Source struct {
	Fsys fs.FS
	Root string
	Name string
}
