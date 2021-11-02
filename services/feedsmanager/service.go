package feedsmanager

import (
	"context"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"time"
)

type FeedsManager struct {
	store FeedStore
}

func New(store FeedStore) *FeedsManager {
	return &FeedsManager{store: store}
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