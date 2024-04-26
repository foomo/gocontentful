package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigFromYAML(t *testing.T) {
	config, err := LoadConfigFromYAML("./sampleconfig.yaml")
	require.NoError(t, err)
	require.Equal(t, "abc123", config.SpaceID)
	require.Equal(t, "dev", config.Environment)
	require.Equal(t, "v1.0.19", config.RequireVersion)
	require.Len(t, config.ContentTypes, 2)
}
