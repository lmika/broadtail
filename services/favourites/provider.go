package favourites

import (
	"context"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
)

type FavouriteStore interface {
	LookupByVideoRef(ctx context.Context, videoRef models.VideoRef) (*models.Favourite, error)
	Save(ctx context.Context, favourite *models.Favourite) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type FeedItemStore interface {
	Get(ctx context.Context, id uuid.UUID) (*models.FeedItem, error)
}
