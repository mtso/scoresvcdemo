package scoresvcdemo

import (
	"time"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
)

// Service middleware (not endpoint)
type Middleware func(Service) Service

type loggingMiddleware struct {
	next Service
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next: next,
			logger: logger,
		}
	}
}

func (mw loggingMiddleware) PostScore(ctx context.Context, s Score) (Score, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "PostScore",
			"id", s.Id,
			"took", time.Since(begin),
			"err", err
		)
	}(time.Now())
	return mw.next.PostScore(ctx, s)
}

func (mw loggingMiddleware) GetScore(ctx context.Context, id string) (Score, error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetScore",
			"id", id,
			"took", time.Since(begin),
			"err", err
		)
	}(time.Now())
	return mw.next.GetScore(ctx, id)
}