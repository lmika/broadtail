package stormstore

import (
	"context"
	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
)

type FeedStore struct {
	db *storm.DB
}

func (js *FeedStore) Delete(ctx context.Context, id uuid.UUID) error {
	return js.db.DeleteStruct(&models.Feed{ID: id})
}

func (js *FeedStore) Get(ctx context.Context, id uuid.UUID) (feed models.Feed, err error) {
	err = js.db.One("ID", id, &feed)
	return feed, err
}

func NewFeedStore(filename string) (*FeedStore, error) {
	db, err := storm.Open(filename)
	if err != nil {
		return nil, err
	}

	return &FeedStore{db: db}, nil
}

func (js *FeedStore) Close() {
	js.db.Close()
}

func (js *FeedStore) List(ctx context.Context) (feeds []models.Feed, err error) {
	err = js.db.Select().OrderBy("Name").Limit(50).Find(&feeds)
	if err == storm.ErrNotFound {
		return []models.Feed{}, nil
	}

	return feeds, err
}

func (js *FeedStore) Save(ctx context.Context, feed *models.Feed) error {
	return js.db.Save(feed)
}
