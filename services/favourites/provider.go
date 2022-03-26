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
	List(ctx context.Context, query models.FeedItemFilter, page int) ([]models.Favourite, error)
}

type VideoMetadata interface {
	GetVideoMetadata(ctx context.Context, youtubeId string) (*models.Video, error)
}

type FeedItemStore interface {
	Get(ctx context.Context, id uuid.UUID) (*models.FeedItem, error)
}

type FeedStore interface {
	Get(ctx context.Context, id uuid.UUID) (models.Feed, error)
}
