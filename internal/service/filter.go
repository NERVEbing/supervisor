package service

import (
	"strings"

	"github.com/NERVEbing/supervisor/internal/domain"
)

type filterService struct {
	rule domain.FilterRule
}

func NewFilterService(rule domain.FilterRule) domain.Filter {
	return &filterService{
		rule: rule,
	}
}

func (f *filterService) ShouldExclude(path string) bool {
	for _, suffix := range f.rule.ExcludeSuffixes {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}

	for _, prefix := range f.rule.ExcludePaths {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}
