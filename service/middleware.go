package service

import (
	"context"
	"time"

	"github.com/benkim0414/geoauth2/client"
	"github.com/go-kit/kit/log"
)

// ServiceMiddleware is a chainable behavior modifier for Service.
type ServiceMiddleware func(Service) Service

func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) PostClient(ctx context.Context, c *client.Client) (client *client.Client, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "post_client", "id", c.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostClient(ctx, c)
}

func (mw loggingMiddleware) GetClient(ctx context.Context, id string) (client *client.Client, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "get_client", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetClient(ctx, id)
}

func (mw loggingMiddleware) DeleteClient(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "delete_client", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteClient(ctx, id)
}
