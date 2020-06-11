package common

import (
	"context"
	"time"
)

// ContextWithDuration creates context with timeout
func ContextWithDuration() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 3000*time.Millisecond)
}
