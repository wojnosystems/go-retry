package core

import (
	"context"
	"time"
)

// Sleep will block until either the ctx expires or the duration is exceeded
func Sleep(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(duration):
	}
	return
}
