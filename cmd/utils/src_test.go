// Copyright 2023 FMNX Linux team.
// This code is covered by GPL license, which can be found in LICENSE file.
// Additional information could be found on official web page: https://fmnx.io/
// Email: help@fmnx.io
package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClean(t *testing.T) {
	ex, err := os.Executable()
	assert.NoError(t, err)
	exPath := filepath.Dir(ex)

	err = os.MkdirAll(exPath+"/test", 0o600)
	assert.NoError(t, err)
	err = os.Chmod(exPath+"/test", 0o777)
	assert.NoError(t, err)

	cleaner := OsHelper{}
	err = cleaner.Clean(exPath)
	assert.NoError(t, err)
	es, err := os.ReadDir(exPath)
	assert.NoError(t, err)

	for _, de := range es {
		assert.True(t, !de.IsDir())
	}
}

func TestExecute(t *testing.T) {
	c := OsHelper{}
	err := c.Execute("ls")
	assert.NoError(t, err)
}
