package git

import (
	"fmt"

	"github.com/NERVEbing/supervisor/internal/domain"

	"github.com/go-git/go-git/v6"
)

type Adapter struct {
	repo *git.Repository
}

func NewAdapter(repoPath string) (domain.Repository, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrRepoNotFound, err)
	}
	return &Adapter{repo: repo}, nil
}
