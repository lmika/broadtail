package feedsmanager

import (
	"context"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/models/ytrss"
	"github.com/pkg/errors"
	"sort"
	"time"
)

type FeedsManager struct {
	store FeedStore
	feedProvider RSSFetcher
}

func New(store FeedStore, feedProvider RSSFetcher) *FeedsManager {
	return &FeedsManager{
		store: store,
		feedProvider: feedProvider,
	}
}

func (fm *FeedsManager) List(ctx context.Context) ([]models.Feed, error) {
	return fm.store.List(ctx)
}

func (fm *FeedsManager) Get(ctx context.Context, id uuid.UUID) (models.Feed, error) {
	return fm.store.Get(ctx, id)
}

func (fm *FeedsManager) Save(ctx context.Context, feed *models.Feed) error {
	if feed.ID == uuid.Nil {
		feed.ID = uuid.New()
		feed.CreatedAt = time.Now()
	}
	return fm.store.Save(ctx, feed)
}

func (fm *FeedsManager) Delete(ctx context.Context, id uuid.UUID) error {
	return fm.store.Delete(ctx, id)
}

func (fm *FeedsManager) RecentFeedItems(ctx context.Context, id uuid.UUID) (entries []ytrss.Entry, err error) {
	feed, err := fm.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	switch feed.Type {
	case models.FeedTypeYoutubeChannel:
		entries, err = fm.feedProvider.GetForChannelID(ctx, feed.ExtID)
	case models.FeedTypeYoutubePlaylist:
		entries, err = fm.feedProvider.GetForPlaylistID(ctx, feed.ExtID)
	default:
		return nil, errors.Errorf("unrecognised feed type: %v", feed.Type)
	}
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Published.After(entries[j].Published)
	})

	return entries, nil
}