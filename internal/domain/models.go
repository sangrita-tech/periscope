package domain

import (
	"io"
)

type Target struct {
	Value            string
	WorkingDirectory string
}

type IgnorePatterns []string

type InspectProjectRequest struct {
	Target Target
	Mode   InspectionMode
	Ignore IgnorePatterns
}

type Repository struct {
	URL      string
	Name     string
	Owner    string
	Branch   string
	SubPath  string
	Provider GitProvider
}

type ProjectScanRequest struct {
	WorkspacePath string
	Ignore        IgnorePatterns
}

type ProjectFile struct {
	Path         string
	RelativePath string
	Content      string
}

type ProjectEntry struct {
	Kind         ProjectEntryKind
	Path         string
	RelativePath string
	Depth        int
	IsLast       bool
}

type Project struct {
	WorkspacePath string
	Files         []ProjectFile
	Entries       []ProjectEntry
}

type InspectionResult struct {
	Text string
}

type CLICommandInput struct {
	Args             []string
	WorkingDirectory string
	Stdout           io.Writer
	Stderr           io.Writer
}

type CLICommandOutput struct {
	Result   InspectionResult
	ExitCode int
	Error    error
}
