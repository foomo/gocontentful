package erm

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testExportFile = "../test/test-space-export.json"

func TestGenerateAPI_EmbedRegionEU(t *testing.T) {
	dir := t.TempDir()
	err := GenerateAPI(context.Background(), dir, "testpkg", "", "", "", testExportFile, nil, "test", "eu")
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(dir, "testpkg", "gocontentfulvolib.go"))
	require.NoError(t, err)

	src := string(content)
	assert.True(t, strings.Contains(src, `.SetRegion("eu")`), "generated lib should call SetRegion(\"eu\")")
}

func TestGenerateAPI_InvalidRegion(t *testing.T) {
	err := GenerateAPI(context.Background(), t.TempDir(), "testpkg", "", "", "", testExportFile, nil, "test", "ap")
	require.Error(t, err)
	assert.ErrorContains(t, err, "unknown region")
}

func TestGenerateAPI_EmbedRegionEmpty(t *testing.T) {
	dir := t.TempDir()
	err := GenerateAPI(context.Background(), dir, "testpkg", "", "", "", testExportFile, nil, "test", "")
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(dir, "testpkg", "gocontentfulvolib.go"))
	require.NoError(t, err)

	src := string(content)
	assert.False(t, strings.Contains(src, `.SetRegion(`), "generated lib should not call SetRegion for default region")
}
