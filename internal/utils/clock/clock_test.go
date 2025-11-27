package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFixedClockNow(t *testing.T) {
	fixedTime := time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)
	clock := FixedClock{FixedTime: fixedTime}

	require.Equal(t, fixedTime, clock.Now())
}

func TestRealClockNow(t *testing.T) {
	rc := RealClock{}
	before := time.Now().Add(-10 * time.Millisecond)
	now := rc.Now()
	after := time.Now().Add(10 * time.Millisecond)

	require.True(t, now.After(before))
	require.True(t, now.Before(after))
}

