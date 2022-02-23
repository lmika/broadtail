package stormstore

import (
	"context"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/pkg/errors"
)

type FeedItemStore struct {
	db *storm.DB
}

func (f *FeedItemStore) PutIfAbsent(ctx context.Context, item *models.FeedItem) error {
	if err := f.db.Select(q.Eq("EntryID", item.EntryID)).First(&models.FeedItem{}); err != nil {
		if !errors.Is(err, storm.ErrNotFound) {
			return err
		}
	} else {
		// Item exists.  Do nothing
		return nil
	}

	item.ID = uuid.New()

	return f.db.Save(item)
}

func NewFeedItemStore(filename string) (*FeedItemStore, error) {
	db, err := storm.Open(filename)
	if err != nil {
		return nil, err
	}

	return &FeedItemStore{db: db}, nil
}

func (f *FeedItemStore) ListRecentsFromAllFeeds(ctx context.Context, limit int) (feedItems []models.FeedItem, err error) {
	err = f.db.Select().OrderBy("Published").Reverse().Limit(limit).Find(&feedItems)
	if err == storm.ErrNotFound {
		return []models.FeedItem{}, nil
	}
	return feedItems, err
}

func (f *FeedItemStore) ListRecent(ctx context.Context, feedID uuid.UUID) (feedItems []models.FeedItem, err error) {
	err = f.db.Select(q.Eq("FeedID", feedID)).OrderBy("Published").Reverse().Limit(50).Find(&feedItems)
	if err == storm.ErrNotFound {
		return []models.FeedItem{}, nil
	}
	return feedItems, err
}

func (f *FeedItemStore) Close() {
	f.db.Close()
}
