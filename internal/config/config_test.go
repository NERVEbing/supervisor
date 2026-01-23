package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromEnv_WithAllEnvVars(t *testing.T) {
	t.Setenv("SUPERVISOR_REPO_PATH", "/test/repo")
	t.Setenv("SUPERVISOR_EXCLUDE_SUFFIXES", ".png,.wasm,.gz")
	t.Setenv("SUPERVISOR_EXCLUDE_PATHS", "vendor/,third_party/")

	cfg := LoadFromEnv()

	assert.Equal(t, "/test/repo", cfg.RepoPath)
	assert.Equal(t, []string{".png", ".wasm", ".gz"}, cfg.ExcludeSuffixes)
	assert.Equal(t, []string{"vendor/", "third_party/"}, cfg.ExcludePaths)
}

func TestLoadFromEnv_WithEmptyEnvVars(t *testing.T) {
	_ = os.Unsetenv("SUPERVISOR_REPO_PATH")
	_ = os.Unsetenv("SUPERVISOR_EXCLUDE_SUFFIXES")
	_ = os.Unsetenv("SUPERVISOR_EXCLUDE_PATHS")

	cfg := LoadFromEnv()

	assert.Equal(t, "", cfg.RepoPath)
	assert.Empty(t, cfg.ExcludeSuffixes)
	assert.Empty(t, cfg.ExcludePaths)
}
