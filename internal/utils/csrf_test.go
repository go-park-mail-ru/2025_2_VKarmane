package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock"
)

func TestGenerateAndValidateCSRF(t *testing.T) {
	fixedClock := clock.FixedClock{
		FixedTime: time.Date(2025, 10, 22, 19, 0, 0, 0, time.Local),
	}

	token, err := GenerateCSRF(fixedClock.FixedTime, "secret")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := ValidateCSRF(token, "secret")
	require.NoError(t, err)
	assert.Equal(t, fixedClock.FixedTime, claims.TimeStamp)

	_, err = ValidateCSRF(token, "wrong")
	assert.Error(t, err)
}
