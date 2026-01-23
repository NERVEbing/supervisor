package config

import (
	"os"
	"strings"
)

// Config holds global configuration from environment variables
type Config struct {
	RepoPath        string
	ExcludeSuffixes []string
	ExcludePaths    []string
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	cfg := &Config{
		RepoPath: os.Getenv("SUPERVISOR_REPO_PATH"),
	}

	if suffixes := os.Getenv("SUPERVISOR_EXCLUDE_SUFFIXES"); suffixes != "" {
		cfg.ExcludeSuffixes = strings.Split(suffixes, ",")
	}

	if paths := os.Getenv("SUPERVISOR_EXCLUDE_PATHS"); paths != "" {
		cfg.ExcludePaths = strings.Split(paths, ",")
	}

	return cfg
}
