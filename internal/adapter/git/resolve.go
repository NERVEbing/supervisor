package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/NERVEbing/supervisor/internal/domain"

	"github.com/go-git/go-git/v6/plumbing"
)

func (a *Adapter) ResolveRef(ctx context.Context, ref string) (string, string, error) {
	hash, err := a.repo.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		return "", "", fmt.Errorf("%w: %s", domain.ErrRefNotFound, ref)
	}

	refType := a.determineRefType(ref)
	return hash.String(), refType, nil
}

func (a *Adapter) determineRefType(ref string) string {
	if strings.HasPrefix(ref, "refs/tags/") || a.isTag(ref) {
		return "tag"
	}

	if strings.HasPrefix(ref, "refs/heads/") || a.isBranch(ref) {
		return "branch"
	}

	return "commit"
}

func (a *Adapter) isTag(ref string) bool {
	tags, err := a.repo.Tags()
	if err != nil {
		return false
	}
	defer tags.Close()

	err = tags.ForEach(func(t *plumbing.Reference) error {
		if t.Name().Short() == ref {
			return domain.ErrStopIteration
		}
		return nil
	})
	return err == domain.ErrStopIteration
}

func (a *Adapter) isBranch(ref string) bool {
	branches, err := a.repo.Branches()
	if err != nil {
		return false
	}
	defer branches.Close()

	err = branches.ForEach(func(b *plumbing.Reference) error {
		if b.Name().Short() == ref {
			return domain.ErrStopIteration
		}
		return nil
	})
	return err == domain.ErrStopIteration
}

func (a *Adapter) CalculateBaseline(ctx context.Context, fromHash, toHash string) (string, string, bool, error) {
	from := plumbing.NewHash(fromHash)
	to := plumbing.NewHash(toHash)

	isAncestor, err := a.isAncestor(from, to)
	if err != nil {
		return "", "", false, fmt.Errorf("failed to check ancestry: %w", err)
	}

	if isAncestor {
		return fromHash, "direct", true, nil
	}

	mergeBase, err := a.findMergeBase(from, to)
	if err != nil {
		return "", "", false, fmt.Errorf("failed to find merge-base: %w", err)
	}

	return mergeBase.String(), "merge-base", false, nil
}

func (a *Adapter) isAncestor(ancestor, descendant plumbing.Hash) (bool, error) {
	ancestorCommit, err := a.repo.CommitObject(ancestor)
	if err != nil {
		return false, err
	}

	descendantCommit, err := a.repo.CommitObject(descendant)
	if err != nil {
		return false, err
	}

	return descendantCommit.IsAncestor(ancestorCommit)
}

func (a *Adapter) findMergeBase(hash1, hash2 plumbing.Hash) (*plumbing.Hash, error) {
	commit1, err := a.repo.CommitObject(hash1)
	if err != nil {
		return nil, err
	}

	commit2, err := a.repo.CommitObject(hash2)
	if err != nil {
		return nil, err
	}

	bases, err := commit1.MergeBase(commit2)
	if err != nil {
		return nil, err
	}

	if len(bases) == 0 {
		return nil, fmt.Errorf("no merge base found")
	}

	hash := bases[0].Hash
	return &hash, nil
}
