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

type RSSFetcher interface {
	GetForChannelID(ctx context.Context, channelID string) ([]ytrss.Entry, error)
	GetForPlaylistID(ctx context.Context, playlistID string) ([]ytrss.Entry, error)
}