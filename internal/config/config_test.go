package config

import (
	"context"
	io "io"
	"testing"

	log "github.com/obalunenko/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func unsetEnv(tb testing.TB) {
	tb.Helper()

	tb.Setenv(portEnv, "")
	tb.Setenv(hostEnv, "")
	tb.Setenv(boxesEnv, "")
	tb.Setenv(levelEnv, "")
	tb.Setenv(formatEnv, "")
}

type noopCloser struct {
	io.Writer
}

func (noopCloser) Close() error {
	return nil
}

func TestLoadDefault(t *testing.T) {
	ctx := context.Background()

	// Disable logging.
	l := log.Init(ctx, log.Params{
		Writer: noopCloser{Writer: io.Discard},
	})

	ctx = log.ContextWithLogger(ctx, l)

	unsetEnv(t)

	t.Run("default", func(t *testing.T) {
		cfg, err := Load(ctx)
		require.NoError(t, err)

		require.Equal(t, DefaultConfig(), cfg)
	})

	t.Run("env", func(t *testing.T) {
		t.Run("port", func(t *testing.T) {
			t.Setenv(portEnv, "8081")

			cfg, err := Load(ctx)
			require.NoError(t, err)

			expected := DefaultConfig()
			expected.HTTP.Port = "8081"

			assert.Equal(t, expected, cfg)
		})
		t.Run("host", func(t *testing.T) {
			t.Setenv(hostEnv, "127.0.0.1")

			cfg, err := Load(ctx)
			require.NoError(t, err)

			expected := DefaultConfig()
			expected.HTTP.Host = "127.0.0.1"

			assert.Equal(t, expected, cfg)
		})
		t.Run("boxes", func(t *testing.T) {
			t.Setenv(boxesEnv, "1,2,3")

			cfg, err := Load(ctx)
			require.NoError(t, err)

			expected := DefaultConfig()
			expected.Pack.Boxes = []uint{1, 2, 3}

			assert.Equal(t, expected, cfg)
		})
		t.Run("boxes - invalid value", func(t *testing.T) {
			t.Setenv(boxesEnv, "sssd212")

			cfg, err := Load(ctx)
			assert.Error(t, err)

			assert.Nil(t, cfg)
		})
		t.Run("level", func(t *testing.T) {
			t.Setenv(levelEnv, "DEBUG")

			cfg, err := Load(ctx)
			require.NoError(t, err)

			expected := DefaultConfig()
			expected.Log.Level = "DEBUG"

			assert.Equal(t, expected, cfg)
		})
		t.Run("format", func(t *testing.T) {
			t.Setenv(formatEnv, "json")

			cfg, err := Load(ctx)
			require.NoError(t, err)

			expected := DefaultConfig()
			expected.Log.Format = "json"

			assert.Equal(t, expected, cfg)
		})
	})
}
