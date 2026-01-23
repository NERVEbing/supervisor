package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildApp_HasDiffCommand(t *testing.T) {
	cmd := buildApp()

	assert.Equal(t, "supervisor", cmd.Name)
	assert.NotEmpty(t, cmd.Commands)

	// Find diff command
	var foundDiff bool
	for _, subCmd := range cmd.Commands {
		if subCmd.Name == "diff" {
			foundDiff = true
			break
		}
	}

	assert.True(t, foundDiff, "should have 'diff' command")
}
