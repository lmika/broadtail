package feedsmanager

import (
	"context"

	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
)

type FeedStore interface {
	List(ctx context.Context) ([]models.Feed, error)
	Save(ctx context.Context, feed *models.Feed) error
	Get(ctx context.Context, id uuid.UUID) (models.Feed, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type FeedItemStore interface {
	Save(ctx context.Context, feedItem *models.FeedItem) error
	Get(ctx context.Context, id uuid.UUID) (*models.FeedItem, error)
	ListRecentsFromAllFeeds(ctx context.Context, filterExpression models.FeedItemFilter, page, count int) ([]models.FeedItem, error)
	ListRecent(ctx context.Context, feedID uuid.UUID, filterExpression models.FeedItemFilter, page int) ([]models.FeedItem, error)
	PutIfAbsent(ctx context.Context, item *models.FeedItem) (wasInserted bool, err error)
}

// type RSSFetcher interface {
// GetForFeed(ctx context.Context, feed models.Feed) ([]ytrss.Entry, error)
// }

type FeedFetcher interface {
	GetForFeed(ctx context.Context, feed models.Feed) ([]models.FetchedFeedItem, error)
	FeedExternalURL(feed models.Feed) (string, error)
	FeedHints(feed models.Feed) models.FeedHints
}

type RulesStore interface {
	List(ctx context.Context) ([]*models.Rule, error)
}

type VideoDownloader interface {
	QueueForDownload(ctx context.Context, videoRef models.VideoRef, feed *models.Feed) error
}
