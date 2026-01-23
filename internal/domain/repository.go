package domain

import "context"

type RefResolver interface {
	ResolveRef(ctx context.Context, ref string) (hash string, refType string, err error)
}

type DiffCalculator interface {
	CalculateDiff(ctx context.Context, fromHash, toHash string) ([]FileChange, DiffStats, error)
}

type HistoryProvider interface {
	GetHistory(ctx context.Context, fromHash, toHash string, excludeMerges bool) ([]Commit, error)
}

type BaselineCalculator interface {
	CalculateBaseline(ctx context.Context, fromHash, toHash string) (baseHash, strategy string, isLinear bool, err error)
}

type MetadataProvider interface {
	GetRepoURL() string
	GetRepoName() string
	GetDiffURL(base, target string) string
}

type Repository interface {
	RefResolver
	DiffCalculator
	HistoryProvider
	BaselineCalculator
	MetadataProvider
}
