package git

import (
	"context"
	"fmt"
	"sort"

	"github.com/NERVEbing/supervisor/internal/domain"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func (a *Adapter) GetHistory(ctx context.Context, fromHash, toHash string, excludeMerges bool) ([]domain.Commit, error) {
	from := plumbing.NewHash(fromHash)
	to := plumbing.NewHash(toHash)

	if from == to {
		return []domain.Commit{}, nil
	}

	cIter, err := a.repo.Log(&git.LogOptions{
		From:  to,
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	var commits []domain.Commit

	err = cIter.ForEach(func(c *object.Commit) error {
		if c.Hash == from {
			return domain.ErrStopIteration
		}

		if excludeMerges && len(c.ParentHashes) > 1 {
			return nil
		}

		commit := domain.Commit{
			Hash:    c.Hash.String(),
			Author:  c.Author.Name,
			Date:    c.Author.When.UTC(),
			Message: c.Message,
			DiffURL: a.buildCommitURL(c.Hash.String()),
		}
		commits = append(commits, commit)
		return nil
	})

	if err != nil && err != domain.ErrStopIteration {
		return nil, err
	}

	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.Before(commits[j].Date)
	})

	return commits, nil
}
