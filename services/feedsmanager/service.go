package feedsmanager

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/lmika/broadtail/services/favourites"

	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/pkg/errors"
)

type FeedsManager struct {
	store         FeedStore
	feedItemStore FeedItemStore
	// rssFeedSource    RSSFetcher
	feedFetcher      FeedFetcher
	favouriteService *favourites.Service
	rulesStore       RulesStore
	videoStore       VideoStore
	videoDownloader  VideoDownloader

	feedUpdateMutex *sync.Mutex
}

func New(
	store FeedStore,
	feedProvider FeedItemStore,
	// rssFeedSource RSSFetcher,
	feedFetcher FeedFetcher,
	favouriteService *favourites.Service,
	rulesStore RulesStore,
	videoStore VideoStore,
	videoDownloader VideoDownloader,
) *FeedsManager {
	return &FeedsManager{
		store:            store,
		feedItemStore:    feedProvider,
		feedFetcher:      feedFetcher,
		rulesStore:       rulesStore,
		favouriteService: favouriteService,
		videoStore:       videoStore,
		videoDownloader:  videoDownloader,
		feedUpdateMutex:  new(sync.Mutex),
	}
}

func (fm *FeedsManager) List(ctx context.Context) ([]models.Feed, error) {
	return fm.store.List(ctx)
}

func (fm *FeedsManager) Get(ctx context.Context, id uuid.UUID) (models.Feed, error) {
	return fm.store.Get(ctx, id)
}

func (fm *FeedsManager) FeedExternalURL(f models.Feed) (string, error) {
	return fm.feedFetcher.FeedExternalURL(f)
}

func (fm *FeedsManager) Save(ctx context.Context, feed *models.Feed) error {
	if feed.ID == uuid.Nil {
		feed.ID = uuid.New()
		feed.CreatedAt = time.Now()

		fs := fm.feedFetcher.FeedHints(*feed)
		feed.Ordering = fs.Ordering
		feed.CheckForUpdates = fs.CheckForUpdates
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

	rules, err := fm.rulesStore.List(ctx)
	if err != nil {
		return errors.Wrapf(err, "cannot get rules")
	}

	return fm.updateFeedItems(ctx, feed, rules)
}

func (fm *FeedsManager) UpdateAllFeeds(ctx context.Context) error {
	allFeeds, err := fm.store.List(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot get all feeds")
	}

	rules, err := fm.rulesStore.List(ctx)
	if err != nil {
		return errors.Wrapf(err, "cannot get rules")
	}

	for _, feed := range allFeeds {
		if !feed.CheckForUpdates {
			continue
		}

		if err := fm.updateFeedItems(ctx, feed, rules); err != nil {
			log.Printf("unable to update feed: %v", err)
		}
	}

	return nil
}

func (fm *FeedsManager) updateFeedItems(ctx context.Context, feed models.Feed, rules []*models.Rule) error {
	fm.feedUpdateMutex.Lock()
	defer fm.feedUpdateMutex.Unlock()

	rssItems, err := fm.feedFetcher.GetForFeed(ctx, feed)
	if err != nil {
		return errors.Wrapf(err, "cannot get feed items from source")
	}

	for _, item := range rssItems {
		feedItem := fm.sourceEntryToFeedItem(&feed, item)

		wasInserted, err := fm.feedItemStore.PutIfAbsent(ctx, &feedItem)
		if err != nil {
			log.Printf("warn: cannot save item %v: %v", feedItem.VideoRef, err)
			continue
		}

		if wasInserted {
			if err := fm.runRulesForFeedItem(ctx, &feed, &feedItem, item, rules); err != nil {
				log.Printf("warn: error running rules for feed item %v", feedItem.ID)
			}
		}
	}

	feed.LastUpdatedAt = time.Now()
	if err := fm.store.Save(ctx, &feed); err != nil {
		return errors.Wrap(err, "cannot update feed")
	}

	return nil
}

func (fm *FeedsManager) runRulesForFeedItem(
	ctx context.Context,
	feed *models.Feed,
	feedItem *models.FeedItem,
	fetchedFeedItem models.FetchedFeedItem,
	rules []*models.Rule,
) error {
	ruleTarget := models.RuleTarget{
		FeedID:      feedItem.FeedID,
		Title:       feedItem.Title,
		Description: fetchedFeedItem.Description,
	}

	// Get all matching rules
	var matchedRules []*models.Rule
	for _, rule := range rules {
		if !rule.Active {
			continue
		}

		if rule.Condition.Matches(ruleTarget) {
			matchedRules = append(matchedRules, rule)
		}
	}
	if len(matchedRules) == 0 {
		return nil
	}

	// Combine the actions
	var combinedAction models.RuleAction
	for _, rule := range matchedRules {
		combinedAction = rule.Action.Combine(combinedAction)
	}

	// Apply the actions
	if combinedAction.Download {
		// Start a download
		if err := fm.videoDownloader.QueueForDownload(ctx, feedItem.VideoRef, feed); err != nil {
			log.Printf("warn: unable to queue download job for feed item %v: %v", feedItem.ID, err)
		}
	}
	if combinedAction.MarkFavourite {
		// Mark as a favourite
		if _, err := fm.favouriteService.FavoriteVideoByOrigin(ctx, models.FavouriteOrigin{
			Type: models.FeedItemOriginType,
			ID:   feedItem.ID.String(),
		}); err != nil {
			log.Printf("warn: unable to add feed item %v as favourite: %v", feedItem.ID, err)
		}
	}
	if combinedAction.MarkDownloaded {
		feedItem.MarkedAsDownloaded = true
		if err := fm.feedItemStore.Save(ctx, feedItem); err != nil {
			log.Printf("warn: unable to update feed item %v to mark as downloaded: %v", feedItem.ID, err)
		}
	}

	return nil
}

func (fm *FeedsManager) sourceEntryToFeedItem(feed *models.Feed, entry models.FetchedFeedItem) models.FeedItem {
	return models.FeedItem{
		VideoRef:  entry.VideoRef,
		FeedID:    feed.ID,
		Title:     entry.Title,
		Link:      entry.Link,
		Published: entry.Published,
	}
}

func (fm *FeedsManager) RecentFeedItems(ctx context.Context, feed *models.Feed, filterExpression models.FeedItemFilter, page int) ([]models.RecentFeedItem, error) {
	feedItems, err := fm.feedItemStore.ListRecent(ctx, feed.ID, filterExpression, page)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get items from feed %v", feed.ID)
	}

	recentFeedItems := make([]models.RecentFeedItem, 0)
	for _, fi := range feedItems {
		var wasDownloaded = fi.MarkedAsDownloaded
		if vs, err := fm.videoStore.FindWithExtID(fi.VideoRef); err == nil && vs != nil {
			wasDownloaded = true
		}

		recentFeedItems = append(recentFeedItems, models.RecentFeedItem{
			Feed:        *feed,
			FeedItem:    fi,
			FavouriteID: fm.favouriteIdForFeedItem(ctx, fi),
			Downloaded:  wasDownloaded,
		})
	}

	return recentFeedItems, nil
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

		var wasDownloaded = fi.MarkedAsDownloaded
		if vs, err := fm.videoStore.FindWithExtID(fi.VideoRef); err == nil && vs != nil {
			wasDownloaded = true
		}

		recentFeedItems = append(recentFeedItems, models.RecentFeedItem{
			Feed:        feed,
			FeedItem:    fi,
			FavouriteID: fm.favouriteIdForFeedItem(ctx, fi),
			Downloaded:  wasDownloaded,
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

func (fm *FeedsManager) favouriteIdForFeedItem(ctx context.Context, feedItem models.FeedItem) string {
	var favouriteId = ""
	f, err := fm.favouriteService.VideoFavourited(ctx, feedItem.VideoRef)
	if err != nil {
		log.Printf("warn: cannot get favourite for item with id: %v", err)
	} else if f != nil {
		favouriteId = f.ID.String()
	}

	return favouriteId
}
