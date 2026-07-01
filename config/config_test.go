package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFromYAML(t *testing.T) {
	t.Run("full config with region", func(t *testing.T) {
		cfg, err := LoadConfigFromYAML("./sampleconfig.yaml")
		require.NoError(t, err)
		require.Equal(t, "abc123", cfg.SpaceID)
		require.Equal(t, "dev", cfg.Environment)
		require.Equal(t, "v1.0.19", cfg.RequireVersion)
		require.Len(t, cfg.ContentTypes, 2)
		assert.Equal(t, "eu", cfg.Region)
	})
	t.Run("region defaults to empty string when absent", func(t *testing.T) {
		f, err := os.CreateTemp("", "gocontentful-test-*.yaml")
		require.NoError(t, err)
		defer os.Remove(f.Name())
		_, err = f.WriteString("spaceId: abc123\n")
		require.NoError(t, err)
		require.NoError(t, f.Close())

		cfg, err := LoadConfigFromYAML(f.Name())
		require.NoError(t, err)
		require.Equal(t, "abc123", cfg.SpaceID)
		assert.Equal(t, "", cfg.Region)
	})
}
