package domain

type FilePath struct {
	Before string
	After  string
}

type Classification struct {
	IsNew       bool
	IsRename    bool
	IsBinary    bool
	IsGenerated bool
	IsTest      bool
	IsConfig    bool
}

type FileHistory struct {
	RelatedCommits []string
}

type FileLineStats struct {
	Added   int
	Deleted int
}

type FileChange struct {
	Path           FilePath
	ChangeType     string
	Language       string
	Lines          FileLineStats
	Classification Classification
	History        FileHistory
}

type DiffSummary struct {
	Files FileStats
	Lines SummaryLineStats
}

type FileStats struct {
	Added    int
	Modified int
	Deleted  int
	Renamed  int
}

type SummaryLineStats struct {
	Added   int
	Deleted int
	Net     int
}

type DiffStats struct {
	BinaryFilesDetected int
	FilesFilteredOut    int
}
