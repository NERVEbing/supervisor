package domain

type Filter interface {
	ShouldExclude(path string) bool
}

type FilterRule struct {
	ExcludeSuffixes []string
	ExcludePaths    []string
}
