package main

import (
	"context"
	"fmt"
	"os"

	"github.com/NERVEbing/supervisor/internal/adapter/git"
	"github.com/NERVEbing/supervisor/internal/adapter/presenter"
	"github.com/NERVEbing/supervisor/internal/config"
	"github.com/NERVEbing/supervisor/internal/domain"
	"github.com/NERVEbing/supervisor/internal/service"

	"github.com/urfave/cli/v3"
)

func runDiff(ctx context.Context, cmd *cli.Command) error {
	cfg := config.LoadFromEnv()

	repoPath := cmd.String("repo")
	if repoPath == "" {
		repoPath = cfg.RepoPath
	}
	if repoPath == "" {
		repoPath = "."
	}

	fromRef := cmd.String("from")
	toRef := cmd.String("to")

	excludeSuffixes := cmd.StringSlice("exclude-suffix")
	if len(excludeSuffixes) == 0 {
		excludeSuffixes = cfg.ExcludeSuffixes
	}

	excludePaths := cmd.StringSlice("exclude-path")
	if len(excludePaths) == 0 {
		excludePaths = cfg.ExcludePaths
	}

	// Dependency Injection
	repo, err := git.NewAdapter(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	filterRule := domain.FilterRule{
		ExcludeSuffixes: excludeSuffixes,
		ExcludePaths:    excludePaths,
	}
	filter := service.NewFilterService(filterRule)

	diffService := service.NewDiffService(repo, filter, filterRule)

	opts := domain.RequestOptions{
		IgnoreMergeCommits: true,
		DetectRenames:      false,
	}

	report, err := diffService.GenerateReport(ctx, fromRef, toRef, opts)
	if err != nil {
		return fmt.Errorf("diff failed: %w", err)
	}

	jsonOutput, err := presenter.ToJSON(report)
	if err != nil {
		return fmt.Errorf("failed to serialize JSON: %w", err)
	}

	_, err = fmt.Fprintln(os.Stdout, string(jsonOutput))
	return err
}
