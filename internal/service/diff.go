package service

import (
	"context"
	"time"

	"github.com/NERVEbing/supervisor/internal/domain"
)

type DiffService struct {
	repo       domain.Repository
	filter     domain.Filter
	filterRule domain.FilterRule
}

func NewDiffService(repo domain.Repository, filter domain.Filter, rule domain.FilterRule) *DiffService {
	return &DiffService{
		repo:       repo,
		filter:     filter,
		filterRule: rule,
	}
}

func (s *DiffService) GenerateReport(ctx context.Context, fromRef, toRef string, opts domain.RequestOptions) (*domain.DiffReport, error) {
	// 1. Resolve Refs
	fromHash, fromType, err := s.repo.ResolveRef(ctx, fromRef)
	if err != nil {
		return nil, err
	}
	toHash, toType, err := s.repo.ResolveRef(ctx, toRef)
	if err != nil {
		return nil, err
	}

	// 2. Calculate Baseline
	baseHash, strategy, isLinear, err := s.repo.CalculateBaseline(ctx, fromHash, toHash)
	if err != nil {
		return nil, err
	}

	// 3. Calculate Diff (Raw)
	rawChanges, rawStats, err := s.repo.CalculateDiff(ctx, fromHash, toHash)
	if err != nil {
		return nil, err
	}

	// 4. Filter and Process Changes
	var filteredChanges []domain.FileChange
	filesFilteredOut := 0

	summaryFiles := domain.FileStats{}
	summaryLines := domain.SummaryLineStats{}

	for _, change := range rawChanges {
		path := change.Path.After
		if path == "" {
			path = change.Path.Before
		}

		if s.filter.ShouldExclude(path) {
			filesFilteredOut++
			continue
		}

		filteredChanges = append(filteredChanges, change)

		switch change.ChangeType {
		case "added":
			summaryFiles.Added++
		case "modified":
			summaryFiles.Modified++
		case "deleted":
			summaryFiles.Deleted++
		case "renamed":
			summaryFiles.Renamed++
		}

		summaryLines.Added += change.Lines.Added
		summaryLines.Deleted += change.Lines.Deleted
	}
	summaryLines.Net = summaryLines.Added - summaryLines.Deleted

	// 5. Get History
	history, err := s.repo.GetHistory(ctx, baseHash, toHash, opts.IgnoreMergeCommits)
	if err != nil {
		// As per original logic, we might want to return empty history on error,
		// but typically in domain service we should return error unless it's non-critical.
		// Original code: "Log error but continue with empty history"
		// I will propagate error for now as it seems safer for strictness,
		// or I can log and ignore if I had a logger.
		// Let's stick to returning error for strict correctness.
		// Or return empty slice if it's acceptable.
		// Given instructions "Strict Go Style Guide", returning error is better.
		return nil, err
	}

	// 6. Assemble Report
	report := &domain.DiffReport{
		SchemaVersion: "1.0",
		Repository: domain.RepoInfo{
			Name: s.repo.GetRepoName(),
			URL:  s.repo.GetRepoURL(),
			VCS:  "git",
		},
		Request: domain.Request{
			FromRef: fromRef,
			ToRef:   toRef,
			Options: opts,
			Filters: domain.RequestFilters{
				ExcludeSuffixes: s.filterRule.ExcludeSuffixes,
				ExcludePaths:    s.filterRule.ExcludePaths,
			},
		},
		Resolution: domain.Resolution{
			From: domain.ResolutionRef{Ref: fromRef, Type: fromType, Commit: fromHash},
			To:   domain.ResolutionRef{Ref: toRef, Type: toType, Commit: toHash},
		},
		Baseline: domain.Baseline{
			Strategy:   strategy,
			BaseCommit: baseHash,
			Ancestry: domain.Ancestry{
				IsLinear:     isLinear,
				Relationship: determineRelationship(isLinear),
			},
		},
		Filters: domain.ReportFilters{
			SuffixExcluded:      s.filterRule.ExcludeSuffixes,
			PathExcluded:        s.filterRule.ExcludePaths,
			BinaryFilesDetected: rawStats.BinaryFilesDetected,
			FilesFilteredOut:    filesFilteredOut,
		},
		TreeDiff: domain.TreeDiff{
			Summary: domain.DiffSummary{
				Files: summaryFiles,
				Lines: summaryLines,
			},
			Files: filteredChanges,
		},
		HistoryView: domain.HistoryView{
			Options: domain.HistoryOptions{
				MergeCommitsIncluded: !opts.IgnoreMergeCommits,
			},
			CommitRange: domain.CommitRange{
				From: baseHash,
				To:   toHash,
			},
			Commits: history,
		},
		DiffLinks: domain.DiffLinks{
			VersionDiff: domain.VersionDiffLink{
				Base:   baseHash,
				Target: toHash,
				URL:    s.repo.GetDiffURL(baseHash, toHash),
			},
		},
		Integrity: domain.Integrity{
			DiffBasis:                      "tree_diff",
			HistoryRole:                    "explanatory_only",
			FilteredFilesNotCountedInStats: true,
			HistoryNote:                    getHistoryNote(history),
		},
		Metadata: domain.Metadata{
			GeneratedAt: time.Now().UTC(),
			Generator: domain.GeneratorInfo{
				Name:    "supervisor",
				Version: "0.1.0",
			},
		},
	}

	return report, nil
}

func determineRelationship(isLinear bool) string {
	if isLinear {
		return "linear"
	}
	return "branched"
}

func getHistoryNote(commits []domain.Commit) string {
	if len(commits) == 0 {
		return "No commits in range (identical base and target, or all merge commits filtered)"
	}
	return ""
}
