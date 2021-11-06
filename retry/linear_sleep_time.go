package retry

import "time"

func linearSleepTime(initial time.Duration, growthFactor float64, iteration uint64) time.Duration {
	return time.Duration(float64(initial) + (float64(initial) * growthFactor * float64(iteration)))
}
