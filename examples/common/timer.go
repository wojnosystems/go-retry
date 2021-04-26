package common

import "time"

func TimeThis(cb func()) time.Duration {
	start := time.Now()
	cb()
	return time.Now().Sub(start)
}

type TimeSet struct {
	start time.Time
	last  time.Time
}

func NewTimeSet() TimeSet {
	now := time.Now()
	return TimeSet{
		start: now,
		last:  now,
	}
}

func (s *TimeSet) SinceLast() time.Duration {
	now := time.Now()
	difference := now.Sub(s.last)
	s.last = now
	return difference
}
