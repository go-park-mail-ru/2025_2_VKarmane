package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	token, err := GenerateJWT(42, "tester", "secret")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := ValidateJWT(token, "secret")
	require.NoError(t, err)
	assert.Equal(t, 42, claims.UserID)
	assert.Equal(t, "tester", claims.Login)

	_, err = ValidateJWT(token, "wrong")
	assert.Error(t, err)
}
