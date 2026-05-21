package domain

type ProjectInspector interface {
	Inspect(InspectProjectRequest) (InspectionResult, error)
}

type TargetResolver interface {
	Resolve(Target) (string, error)
}

type RepositoryFetcher interface {
	Fetch(Repository) (string, error)
}

type ProjectScanner interface {
	Scan(ProjectScanRequest) (Project, error)
}

type ProjectRenderer interface {
	Render(Project) (InspectionResult, error)
}

type ClipboardWriter interface {
	Write(InspectionResult) error
}

type CLICommand interface {
	Run(CLICommandInput) CLICommandOutput
}
