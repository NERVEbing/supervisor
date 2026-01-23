package git

import (
	"context"
	"fmt"

	"github.com/NERVEbing/supervisor/internal/domain"

	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func (a *Adapter) CalculateDiff(ctx context.Context, fromHash, toHash string) ([]domain.FileChange, domain.DiffStats, error) {
	from := plumbing.NewHash(fromHash)
	to := plumbing.NewHash(toHash)

	if from == to {
		return []domain.FileChange{}, domain.DiffStats{}, nil
	}

	fromCommit, err := a.repo.CommitObject(from)
	if err != nil {
		return nil, domain.DiffStats{}, fmt.Errorf("failed to get from commit: %w", err)
	}

	toCommit, err := a.repo.CommitObject(to)
	if err != nil {
		return nil, domain.DiffStats{}, fmt.Errorf("failed to get to commit: %w", err)
	}

	fromTree, err := fromCommit.Tree()
	if err != nil {
		return nil, domain.DiffStats{}, fmt.Errorf("failed to get from tree: %w", err)
	}

	toTree, err := toCommit.Tree()
	if err != nil {
		return nil, domain.DiffStats{}, fmt.Errorf("failed to get to tree: %w", err)
	}

	changes, err := fromTree.Diff(toTree)
	if err != nil {
		return nil, domain.DiffStats{}, fmt.Errorf("failed to calculate tree diff: %w", err)
	}

	return a.convertChanges(changes, fromTree, toTree)
}

func (a *Adapter) convertChanges(changes object.Changes, fromTree, toTree *object.Tree) ([]domain.FileChange, domain.DiffStats, error) {
	var result []domain.FileChange
	var stats domain.DiffStats

	for _, change := range changes {
		path := getChangePath(change)
		fc, isBinary, err := a.convertSingleChange(change, fromTree, toTree)
		if err != nil {
			return nil, stats, fmt.Errorf("failed to convert change for %s: %w", path, err)
		}

		if isBinary {
			stats.BinaryFilesDetected++
		}

		result = append(result, fc)
	}

	return result, stats, nil
}

func (a *Adapter) convertSingleChange(change *object.Change, fromTree, toTree *object.Tree) (domain.FileChange, bool, error) {
	var changeType string
	var pathBefore, pathAfter string

	switch {
	case change.From.Name == "" && change.To.Name != "":
		changeType = "added"
		pathAfter = change.To.Name
	case change.From.Name != "" && change.To.Name == "":
		changeType = "deleted"
		pathBefore = change.From.Name
	case change.From.Name != change.To.Name:
		changeType = "renamed"
		pathBefore = change.From.Name
		pathAfter = change.To.Name
	default:
		changeType = "modified"
		pathBefore = change.From.Name
		pathAfter = change.To.Name
	}

	isBinary, err := isBinaryFile(change, fromTree, toTree)
	if err != nil {
		return domain.FileChange{}, false, err
	}

	var lineStats domain.FileLineStats
	if !isBinary {
		stats, err := a.calculateLineStats(change)
		if err != nil {
			return domain.FileChange{}, isBinary, err
		}
		lineStats = stats
	}

	path := getChangePath(change)
	language := detectLanguage(path)

	classification := domain.Classification{
		IsNew:       changeType == "added",
		IsRename:    changeType == "renamed",
		IsBinary:    isBinary,
		IsGenerated: isGeneratedFile(path),
		IsTest:      isTestFile(path),
		IsConfig:    isConfigFile(path),
	}

	return domain.FileChange{
		Path: domain.FilePath{
			Before: pathBefore,
			After:  pathAfter,
		},
		ChangeType:     changeType,
		Language:       language,
		Lines:          lineStats,
		Classification: classification,
		History: domain.FileHistory{
			RelatedCommits: []string{},
		},
	}, isBinary, nil
}

func (a *Adapter) calculateLineStats(change *object.Change) (domain.FileLineStats, error) {
	patch, err := change.Patch()
	if err != nil {
		return domain.FileLineStats{}, err
	}

	stats := patch.Stats()
	var added, deleted int
	for _, fileStat := range stats {
		added += fileStat.Addition
		deleted += fileStat.Deletion
	}

	return domain.FileLineStats{
		Added:   added,
		Deleted: deleted,
	}, nil
}
