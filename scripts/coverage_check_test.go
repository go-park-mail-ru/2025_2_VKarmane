package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseProfile(t *testing.T) {
	content := "mode: atomic\nexample/file.go:1.1,3.0 2 1\nexample/file.go:4.0,5.0 3 0\n"
	tmp := t.TempDir() + "/cover.out"
	require.NoError(t, os.WriteFile(tmp, []byte(content), 0o600))

	total, covered, err := parseProfile(tmp)
	require.NoError(t, err)
	require.Equal(t, 5.0, total)
	require.Equal(t, 2.0, covered)
}



