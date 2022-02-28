package feedsmanager

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/models/ytrss"
	"github.com/pkg/errors"
)

type FeedsManager struct {
	store         FeedStore
	feedItemStore FeedItemStore
	rssFeedSource RSSFetcher

	feedUpdateMutex *sync.Mutex
}

func New(store FeedStore, feedProvider FeedItemStore, rssFeedSource RSSFetcher) *FeedsManager {
	return &FeedsManager{
		store:           store,
		feedItemStore:   feedProvider,
		rssFeedSource:   rssFeedSource,
		feedUpdateMutex: new(sync.Mutex),
	}
}

func (fm *FeedsManager) List(ctx context.Context) ([]models.Feed, error) {
	return fm.store.List(ctx)
}

func (fm *FeedsManager) Get(ctx context.Context, id uuid.UUID) (models.Feed, error) {
	return fm.store.Get(ctx, id)
}

func (fm *FeedsManager) FeedExternalURL(f models.Feed) (string, error) {
	switch f.Type {
	case models.FeedTypeYoutubeChannel:
		return fmt.Sprintf("https://www.youtube.com/channel/%v", f.ExtID), nil
	case models.FeedTypeYoutubePlaylist:
		return fmt.Sprintf("https://www.youtube.com/playlist/%v", f.ExtID), nil
	}
	return "", errors.Errorf("external url unsupported for feed type: %v", f.Type)
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

	return fm.updateFeedItems(ctx, feed)
}

func (fm *FeedsManager) UpdateAllFeeds(ctx context.Context) error {
	allFeeds, err := fm.store.List(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot get all feeds")
	}

	for _, feed := range allFeeds {
		if err := fm.updateFeedItems(ctx, feed); err != nil {
			log.Printf("unable to update feed: %v", err)
		}
	}

	return nil
}

func (fm *FeedsManager) updateFeedItems(ctx context.Context, feed models.Feed) error {
	fm.feedUpdateMutex.Lock()
	defer fm.feedUpdateMutex.Unlock()

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

	feed.LastUpdatedAt = time.Now()
	if err := fm.store.Save(ctx, &feed); err != nil {
		return errors.Wrap(err, "cannot update feed")
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

func (fm *FeedsManager) RecentFeedItems(ctx context.Context, id uuid.UUID, filterExpression models.FeedItemFilter, page int) (entries []models.FeedItem, err error) {
	return fm.feedItemStore.ListRecent(ctx, id, filterExpression, page)
}

func (fm *FeedsManager) RecentFeedItemsFromAllFeeds(ctx context.Context, filterExpression models.FeedItemFilter, page, count int) ([]models.RecentFeedItem, error) {
	feedItems, err := fm.feedItemStore.ListRecentsFromAllFeeds(ctx, filterExpression, page, count)
	if err != nil {
		return nil, errors.Wrap(err, "cannot list recent feed items")
	}

	recentFeedItems := make([]models.RecentFeedItem, 0)
	for _, fi := range feedItems {
		feed, err := fm.store.Get(ctx, fi.FeedID)
		if err != nil {
			log.Printf("warn: cannot get feed with id: %v", err)
		}

		recentFeedItems = append(recentFeedItems, models.RecentFeedItem{
			Feed:     feed,
			FeedItem: fi,
		})
	}

	return recentFeedItems, nil
}

func (fm *FeedsManager) GetFeedItem(ctx context.Context, feedItemID uuid.UUID) (*models.FeedItem, error) {
	return fm.feedItemStore.Get(ctx, feedItemID)
}

func (fm *FeedsManager) SaveFeedItem(ctx context.Context, feedItem *models.FeedItem) error {
	return fm.feedItemStore.Save(ctx, feedItem)
}
