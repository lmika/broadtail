package feedsmanager

import (
	"context"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/models/ytrss"
)

type FeedStore interface {
	List(ctx context.Context) ([]models.Feed, error)
	Save(ctx context.Context, feed *models.Feed) error
	Get(ctx context.Context, id uuid.UUID) (models.Feed, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type FeedItemStore interface {
	ListRecent(ctx context.Context, feedID uuid.UUID) ([]models.FeedItem, error)
	PutIfAbsent(ctx context.Context, item *models.FeedItem) error
}

type RSSFetcher interface {
	GetForFeed(ctx context.Context, feed models.Feed) ([]ytrss.Entry, error)
}
