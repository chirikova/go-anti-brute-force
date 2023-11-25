package logger

import (
	"bytes"
	"testing"

	"github.com/chirikova/go-anti-brute-force/internal/config"
	"github.com/stretchr/testify/require"
)

var testMessage = "test log message"

func TestLogger(t *testing.T) {
	cfg, _ := config.InitConfig("../../configs/config.yaml")
	buffer := &bytes.Buffer{}
	logger, err := New(cfg.Logger, buffer)

	t.Run("successful create logger", func(t *testing.T) {
		require.NoError(t, err)
	})

	t.Run("successful info", func(t *testing.T) {
		buffer.Reset()
		logger.Info(testMessage)
		require.Contains(t, buffer.String(), testMessage)
		require.Contains(t, buffer.String(), "level=info")
	})

	t.Run("successful warning", func(t *testing.T) {
		buffer.Reset()
		logger.Warn(testMessage)
		require.Contains(t, buffer.String(), testMessage)
		require.Contains(t, buffer.String(), "level=warn")
	})

	t.Run("successful error", func(t *testing.T) {
		buffer.Reset()
		logger.Error(testMessage)
		require.Contains(t, buffer.String(), testMessage)
		require.Contains(t, buffer.String(), "level=error")
	})

	t.Run("no debug with logger level info", func(t *testing.T) {
		buffer.Reset()
		logger.Debug(testMessage)
		require.NotContains(t, buffer.String(), testMessage)
		require.NotContains(t, buffer.String(), "DEBUG")
	})

	t.Run("has debug with logger level debug", func(t *testing.T) {
		cfg.Logger.Level = "DEBUG"
		buffer.Reset()
		logger, _ = New(cfg.Logger, buffer)
		logger.Debug(testMessage)
		require.Contains(t, buffer.String(), testMessage)
		require.Contains(t, buffer.String(), "level=debug")
	})
}
