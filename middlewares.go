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

func makeLoggingFunc(logger log.Logger, methodName string, err error) func(time.Time) {
	return func(begin time.Time) {
		logger.Log("method", methodName, "took", time.Since(begin), "err", err)
	}
}

func (mw loggingMiddleware) GetHealth(ctx context.Context) (health Health, err error) {
	defer makeLoggingFunc(mw.logger, "GetHealth", err)
	return mw.next.GetHealth(ctx)
}

func (mw loggingMiddleware) GetZerglings(ctx context.Context) (zerglings []Zergling, err error) {
	defer makeLoggingFunc(mw.logger, "GetZerglings", err)
	return mw.next.GetZerglings(ctx)
}

func (mw loggingMiddleware) PostZerglings(ctx context.Context) (zergling Zergling, err error) {
	defer makeLoggingFunc(mw.logger, "PostZerglings", err)
	return mw.next.PostZerglings(ctx)
}

func (mw loggingMiddleware) GetZerglingByID(ctx context.Context, id string) (zergling Zergling, err error) {
	defer makeLoggingFunc(mw.logger, "PostZerglings", err)
	return mw.next.GetZerglingByID(ctx, id)
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
