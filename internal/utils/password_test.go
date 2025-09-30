package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashAndVerifyPassword(t *testing.T) {
	password := "SuperSecret123!"

	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	ok, err := VerifyPassword(password, hash)
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = VerifyPassword("wrong-pass", hash)
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestVerifyPasswordInvalidHash(t *testing.T) {
	_, err := VerifyPassword("pass", "invalid-hash-format")
	assert.Error(t, err)
}
