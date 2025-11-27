package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithLoggerAndFromContext(t *testing.T) {
	l := NewSlogLogger()
	ctx := context.Background()

	require.Nil(t, FromContext(ctx))

	ctx = WithLogger(ctx, l)
	require.Equal(t, l, FromContext(ctx))
}

func TestNewSlogLoggerWithLevel(t *testing.T) {
	levels := []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError, 123}
	for _, lvl := range levels {
		t.Run(lvl.String(), func(t *testing.T) {
			require.NotNil(t, NewSlogLoggerWithLevel(lvl))
		})
	}
}

func TestNewSlogLoggerWithFile_InvalidPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing-dir", "app.log")
	l, err := NewSlogLoggerWithFile(path, slog.LevelInfo)
	require.Error(t, err)
	require.Nil(t, l)
}

func TestNewSlogLoggerWithFile_Success(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")
	l, err := NewSlogLoggerWithFile(path, slog.LevelInfo)
	require.NoError(t, err)
	require.NotNil(t, l)

	slogLogger, ok := l.(*SlogLogger)
	require.True(t, ok)
	t.Cleanup(func() { require.NoError(t, slogLogger.Close()) })
	slogLogger.Info("hello")

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.True(t, len(data) > 0)
}

func TestNewSlogLoggerWithFileAndConsole(t *testing.T) {
	path := filepath.Join(t.TempDir(), "combined.log")
	l, err := NewSlogLoggerWithFileAndConsole(path, slog.LevelWarn)
	require.NoError(t, err)
	require.NotNil(t, l)

	slogLogger, ok := l.(*SlogLogger)
	require.True(t, ok)
	t.Cleanup(func() { require.NoError(t, slogLogger.Close()) })
	slogLogger.Warn("combined")

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.True(t, strings.Contains(string(data), "combined"))
}

func TestMultiWriter_Write(t *testing.T) {
	file := newTempFile(t)
	console := newTempFile(t)
	mw := &MultiWriter{file: file, console: console}

	n, err := mw.Write([]byte("data"))
	require.NoError(t, err)
	require.Equal(t, len("data"), n)

	require.Equal(t, "data", readAll(t, file))
	require.Equal(t, "data", readAll(t, console))
}

func newTempFile(t *testing.T) *os.File {
	t.Helper()
	dir := t.TempDir()
	f, err := os.CreateTemp(dir, "*.log")
	require.NoError(t, err)
	t.Cleanup(func() { _ = f.Close() })
	return f
}

func readAll(t *testing.T, f *os.File) string {
	t.Helper()
	_, err := f.Seek(0, io.SeekStart)
	require.NoError(t, err)
	data, err := io.ReadAll(f)
	require.NoError(t, err)
	return string(data)
}
