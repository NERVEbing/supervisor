package domain

import "time"

type DiffReport struct {
	SchemaVersion string
	Repository    RepoInfo
	Request       Request
	Resolution    Resolution
	Baseline      Baseline
	Filters       ReportFilters
	TreeDiff      TreeDiff
	HistoryView   HistoryView
	DiffLinks     DiffLinks
	Integrity     Integrity
	Metadata      Metadata
}

type RepoInfo struct {
	Name string
	URL  string
	VCS  string
}

type Request struct {
	FromRef string
	ToRef   string
	Options RequestOptions
	Filters RequestFilters
}

type RequestOptions struct {
	IgnoreMergeCommits bool
	DetectRenames      bool
}

type RequestFilters struct {
	ExcludeSuffixes []string
	ExcludePaths    []string
}

type Resolution struct {
	From ResolutionRef
	To   ResolutionRef
}

type ResolutionRef struct {
	Ref    string
	Type   string
	Commit string
}

type Baseline struct {
	Strategy   string
	BaseCommit string
	Ancestry   Ancestry
}

type Ancestry struct {
	IsLinear     bool
	Relationship string
}

type ReportFilters struct {
	SuffixExcluded      []string
	PathExcluded        []string
	BinaryFilesDetected int
	FilesFilteredOut    int
}

type TreeDiff struct {
	Summary DiffSummary
	Files   []FileChange
}

type HistoryView struct {
	Options     HistoryOptions
	CommitRange CommitRange
	Commits     []Commit
}

type HistoryOptions struct {
	MergeCommitsIncluded bool
}

type CommitRange struct {
	From string
	To   string
}

type Commit struct {
	Hash    string
	Author  string
	Date    time.Time
	Message string
	DiffURL string
}

type DiffLinks struct {
	VersionDiff VersionDiffLink
}

type VersionDiffLink struct {
	Base   string
	Target string
	URL    string
}

type Integrity struct {
	DiffBasis                      string
	HistoryRole                    string
	FilteredFilesNotCountedInStats bool
	HistoryNote                    string
}

type Metadata struct {
	GeneratedAt time.Time
	Generator   GeneratorInfo
}

type GeneratorInfo struct {
	Name    string
	Version string
}
