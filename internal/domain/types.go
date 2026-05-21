package domain

type InspectionMode int

const (
	InspectionModeContent InspectionMode = iota
	InspectionModeTree
)

type ProjectEntryKind int

const (
	ProjectEntryKindDirectory ProjectEntryKind = iota
	ProjectEntryKindFile
)

type GitProvider string

const (
	GitProviderGitHub  GitProvider = "GITHUB"
	GitProviderGitLab  GitProvider = "GITLAB"
	GitProviderUnknown GitProvider = "UNKNOWN"
)
