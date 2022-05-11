package feedfetchers

import (
	"context"

	"github.com/lmika/broadtail/models"
)

type FeedDriver interface {
	GetForFeed(ctx context.Context, feed models.Feed) ([]models.FetchedFeedItem, error)
	FeedExternalURL(feed models.Feed) (string, error)
}
