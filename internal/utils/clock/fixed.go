package clock

import "time"

type FixedClock struct {
	FixedTime time.Time
}

func (f FixedClock) Now() time.Time {
	return f.FixedTime
}
