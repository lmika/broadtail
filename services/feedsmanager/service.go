package feedsmanager

import (
	"context"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/models/ytrss"
	"github.com/pkg/errors"
	"log"
	"time"
)

type FeedsManager struct {
	store         FeedStore
	feedItemStore FeedItemStore
	rssFeedSource RSSFetcher
}

func New(store FeedStore, feedProvider FeedItemStore, rssFeedSource RSSFetcher) *FeedsManager {
	return &FeedsManager{
		store:         store,
		feedItemStore: feedProvider,
		rssFeedSource: rssFeedSource,
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

func (fm *FeedsManager) UpdateFeed(ctx context.Context, id uuid.UUID) error {
	feed, err := fm.Get(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "cannot get feed")
	}

	rssItems, err := fm.rssFeedSource.GetForFeed(ctx, feed)
	if err != nil {
		return errors.Wrapf(err, "cannot get feed items from source")
	}

	for _, item := range rssItems {
		feedItem := fm.sourceEntryToFeedItem(&feed, item)
		if err := fm.feedItemStore.PutIfAbsent(ctx, &feedItem); err != nil {
			log.Printf("warn: cannot save item %v: %v", feedItem.VideoID, err)
		}
	}

	return nil
}

func (fm *FeedsManager) sourceEntryToFeedItem(feed *models.Feed, entry ytrss.Entry) models.FeedItem {
	return models.FeedItem{
		FeedID:    feed.ID,
		EntryID:   entry.VideoID,
		Title:     entry.Title,
		Link:      entry.Link,
		Published: entry.Published,
	}
}

func (fm *FeedsManager) RecentFeedItems(ctx context.Context, id uuid.UUID) (entries []models.FeedItem, err error) {
	return fm.feedItemStore.ListRecent(ctx, id)
	/*
		feed, err := fm.store.Get(ctx, id)
		if err != nil {
			return nil, err
		}

		switch feed.Type {
		case models.FeedTypeYoutubeChannel:
			entries, err = fm.feedItemStore.GetForChannelID(ctx, feed.ExtID)
		case models.FeedTypeYoutubePlaylist:
			entries, err = fm.feedItemStore.GetForPlaylistID(ctx, feed.ExtID)
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
	*/
}
