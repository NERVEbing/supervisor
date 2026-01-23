package presenter

import (
	"encoding/json"
	"time"

	"github.com/NERVEbing/supervisor/internal/domain"
)

type jsonDiffReport struct {
	SchemaVersion string          `json:"schema_version"`
	Repository    jsonRepository  `json:"repository"`
	Request       jsonRequest     `json:"request"`
	Resolution    jsonResolution  `json:"resolution"`
	Baseline      jsonBaseline    `json:"baseline"`
	Filters       jsonFilters     `json:"filters"`
	TreeDiff      jsonTreeDiff    `json:"tree_diff"`
	HistoryView   jsonHistoryView `json:"history_view"`
	DiffLinks     jsonDiffLinks   `json:"diff_links"`
	Integrity     jsonIntegrity   `json:"integrity"`
	Metadata      jsonMetadata    `json:"metadata"`
}

type jsonRepository struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	VCS  string `json:"vcs"`
}

type jsonRequest struct {
	FromRef string             `json:"from_ref"`
	ToRef   string             `json:"to_ref"`
	Options jsonRequestOptions `json:"options"`
	Filters jsonRequestFilters `json:"filters"`
}

type jsonRequestOptions struct {
	IgnoreMergeCommits bool `json:"ignore_merge_commits"`
	DetectRenames      bool `json:"detect_renames"`
}

type jsonRequestFilters struct {
	ExcludeSuffixes []string `json:"exclude_suffixes"`
	ExcludePaths    []string `json:"exclude_paths"`
}

type jsonResolution struct {
	From jsonResolutionRef `json:"from"`
	To   jsonResolutionRef `json:"to"`
}

type jsonResolutionRef struct {
	Ref    string `json:"ref"`
	Type   string `json:"type"`
	Commit string `json:"commit"`
}

type jsonBaseline struct {
	Strategy   string       `json:"strategy"`
	BaseCommit string       `json:"base_commit"`
	Ancestry   jsonAncestry `json:"ancestry"`
}

type jsonAncestry struct {
	IsLinear     bool   `json:"is_linear"`
	Relationship string `json:"relationship"`
}

type jsonFilters struct {
	SuffixExcluded      []string `json:"suffix_excluded"`
	PathExcluded        []string `json:"path_excluded"`
	BinaryFilesDetected int      `json:"binary_files_detected"`
	FilesFilteredOut    int      `json:"files_filtered_out"`
}

type jsonTreeDiff struct {
	Summary jsonDiffSummary  `json:"summary"`
	Files   []jsonFileChange `json:"files"`
}

type jsonDiffSummary struct {
	Files jsonFileStats        `json:"files"`
	Lines jsonSummaryLineStats `json:"lines"`
}

type jsonFileStats struct {
	Added    int `json:"added"`
	Modified int `json:"modified"`
	Deleted  int `json:"deleted"`
	Renamed  int `json:"renamed"`
}

type jsonSummaryLineStats struct {
	Added   int `json:"added"`
	Deleted int `json:"deleted"`
	Net     int `json:"net"`
}

type jsonFileChange struct {
	Path           jsonFilePath       `json:"path"`
	ChangeType     string             `json:"change_type"`
	Language       string             `json:"language"`
	Lines          jsonFileLineStats  `json:"lines"`
	Classification jsonClassification `json:"classification"`
	History        jsonFileHistory    `json:"history"`
}

type jsonFilePath struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

type jsonFileLineStats struct {
	Added   int `json:"added"`
	Deleted int `json:"deleted"`
}

type jsonClassification struct {
	IsNew       bool `json:"is_new"`
	IsRename    bool `json:"is_rename"`
	IsBinary    bool `json:"is_binary"`
	IsGenerated bool `json:"is_generated"`
	IsTest      bool `json:"is_test"`
	IsConfig    bool `json:"is_config"`
}

type jsonFileHistory struct {
	RelatedCommits []string `json:"related_commits"`
}

type jsonHistoryView struct {
	Options     jsonHistoryOptions `json:"options"`
	CommitRange jsonCommitRange    `json:"commit_range"`
	Commits     []jsonCommit       `json:"commits"`
}

type jsonHistoryOptions struct {
	MergeCommitsIncluded bool `json:"merge_commits_included"`
}

type jsonCommitRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type jsonCommit struct {
	Hash    string    `json:"hash"`
	Author  string    `json:"author"`
	Date    time.Time `json:"date"`
	Message string    `json:"message"`
	DiffURL string    `json:"diff_url"`
}

type jsonDiffLinks struct {
	VersionDiff jsonVersionDiffLink `json:"version_diff"`
}

type jsonVersionDiffLink struct {
	Base   string `json:"base"`
	Target string `json:"target"`
	URL    string `json:"url"`
}

type jsonIntegrity struct {
	DiffBasis                      string `json:"diff_basis"`
	HistoryRole                    string `json:"history_role"`
	FilteredFilesNotCountedInStats bool   `json:"filtered_files_not_counted_in_stats"`
	HistoryNote                    string `json:"history_note"`
}

type jsonMetadata struct {
	GeneratedAt time.Time         `json:"generated_at"`
	Generator   jsonGeneratorInfo `json:"generator"`
}

type jsonGeneratorInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func ToJSON(r *domain.DiffReport) ([]byte, error) {
	dto := mapToDTO(r)
	return json.MarshalIndent(dto, "", "  ")
}

func mapToDTO(r *domain.DiffReport) jsonDiffReport {
	files := make([]jsonFileChange, len(r.TreeDiff.Files))
	for i, f := range r.TreeDiff.Files {
		files[i] = jsonFileChange{
			Path: jsonFilePath{
				Before: f.Path.Before,
				After:  f.Path.After,
			},
			ChangeType: f.ChangeType,
			Language:   f.Language,
			Lines: jsonFileLineStats{
				Added:   f.Lines.Added,
				Deleted: f.Lines.Deleted,
			},
			Classification: jsonClassification{
				IsNew:       f.Classification.IsNew,
				IsRename:    f.Classification.IsRename,
				IsBinary:    f.Classification.IsBinary,
				IsGenerated: f.Classification.IsGenerated,
				IsTest:      f.Classification.IsTest,
				IsConfig:    f.Classification.IsConfig,
			},
			History: jsonFileHistory{
				RelatedCommits: f.History.RelatedCommits,
			},
		}
	}

	commits := make([]jsonCommit, len(r.HistoryView.Commits))
	for i, c := range r.HistoryView.Commits {
		commits[i] = jsonCommit{
			Hash:    c.Hash,
			Author:  c.Author,
			Date:    c.Date,
			Message: c.Message,
			DiffURL: c.DiffURL,
		}
	}

	return jsonDiffReport{
		SchemaVersion: r.SchemaVersion,
		Repository: jsonRepository{
			Name: r.Repository.Name,
			URL:  r.Repository.URL,
			VCS:  r.Repository.VCS,
		},
		Request: jsonRequest{
			FromRef: r.Request.FromRef,
			ToRef:   r.Request.ToRef,
			Options: jsonRequestOptions{
				IgnoreMergeCommits: r.Request.Options.IgnoreMergeCommits,
				DetectRenames:      r.Request.Options.DetectRenames,
			},
			Filters: jsonRequestFilters{
				ExcludeSuffixes: r.Request.Filters.ExcludeSuffixes,
				ExcludePaths:    r.Request.Filters.ExcludePaths,
			},
		},
		Resolution: jsonResolution{
			From: jsonResolutionRef{
				Ref:    r.Resolution.From.Ref,
				Type:   r.Resolution.From.Type,
				Commit: r.Resolution.From.Commit,
			},
			To: jsonResolutionRef{
				Ref:    r.Resolution.To.Ref,
				Type:   r.Resolution.To.Type,
				Commit: r.Resolution.To.Commit,
			},
		},
		Baseline: jsonBaseline{
			Strategy:   r.Baseline.Strategy,
			BaseCommit: r.Baseline.BaseCommit,
			Ancestry: jsonAncestry{
				IsLinear:     r.Baseline.Ancestry.IsLinear,
				Relationship: r.Baseline.Ancestry.Relationship,
			},
		},
		Filters: jsonFilters{
			SuffixExcluded:      r.Filters.SuffixExcluded,
			PathExcluded:        r.Filters.PathExcluded,
			BinaryFilesDetected: r.Filters.BinaryFilesDetected,
			FilesFilteredOut:    r.Filters.FilesFilteredOut,
		},
		TreeDiff: jsonTreeDiff{
			Summary: jsonDiffSummary{
				Files: jsonFileStats{
					Added:    r.TreeDiff.Summary.Files.Added,
					Modified: r.TreeDiff.Summary.Files.Modified,
					Deleted:  r.TreeDiff.Summary.Files.Deleted,
					Renamed:  r.TreeDiff.Summary.Files.Renamed,
				},
				Lines: jsonSummaryLineStats{
					Added:   r.TreeDiff.Summary.Lines.Added,
					Deleted: r.TreeDiff.Summary.Lines.Deleted,
					Net:     r.TreeDiff.Summary.Lines.Net,
				},
			},
			Files: files,
		},
		HistoryView: jsonHistoryView{
			Options: jsonHistoryOptions{
				MergeCommitsIncluded: r.HistoryView.Options.MergeCommitsIncluded,
			},
			CommitRange: jsonCommitRange{
				From: r.HistoryView.CommitRange.From,
				To:   r.HistoryView.CommitRange.To,
			},
			Commits: commits,
		},
		DiffLinks: jsonDiffLinks{
			VersionDiff: jsonVersionDiffLink{
				Base:   r.DiffLinks.VersionDiff.Base,
				Target: r.DiffLinks.VersionDiff.Target,
				URL:    r.DiffLinks.VersionDiff.URL,
			},
		},
		Integrity: jsonIntegrity{
			DiffBasis:                      r.Integrity.DiffBasis,
			HistoryRole:                    r.Integrity.HistoryRole,
			FilteredFilesNotCountedInStats: r.Integrity.FilteredFilesNotCountedInStats,
			HistoryNote:                    r.Integrity.HistoryNote,
		},
		Metadata: jsonMetadata{
			GeneratedAt: r.Metadata.GeneratedAt,
			Generator: jsonGeneratorInfo{
				Name:    r.Metadata.Generator.Name,
				Version: r.Metadata.Generator.Version,
			},
		},
	}
}
