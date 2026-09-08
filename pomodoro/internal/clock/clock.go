package clock

import "time"

type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now()
}

type FixedClock struct {
	Current time.Time
}

func (clock FixedClock) Now() time.Time {
	return clock.Current
}
