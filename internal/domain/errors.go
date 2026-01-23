package domain

import "errors"

var (
	ErrRefNotFound   = errors.New("reference not found")
	ErrRepoNotFound  = errors.New("repository not found")
	ErrStopIteration = errors.New("stop iteration")
)
