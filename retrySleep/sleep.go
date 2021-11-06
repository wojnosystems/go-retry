package retrySleep

import (
	"context"
	"time"
)

// WithContext will block until either the ctx expires or the duration is exceeded
func WithContext(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(duration):
	}
	return
}
