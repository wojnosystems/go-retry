package retry

import (
	"math"
	"time"
)

func exponentialSleepTime(initial time.Duration, growthFactor float64, iteration uint64) time.Duration {
	return time.Duration(
		float64(initial) * math.Pow(1.0+growthFactor, float64(iteration)))
}
