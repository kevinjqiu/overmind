package overmind

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware
type Middleware func(Service) Service

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) GetHealth(ctx context.Context) (health Health, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetHealth", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetHealth(ctx)
}

// LoggingMiddleware provides logging for the service
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}
