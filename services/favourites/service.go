package favourites

import (
	"context"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/pkg/errors"
	"time"
)

type Service struct {
	store         FavouriteStore
	feedItemStore FeedItemStore
}

func NewService(store FavouriteStore, feedItemStore FeedItemStore) *Service {
	return &Service{
		store:         store,
		feedItemStore: feedItemStore,
	}
}

// VideoFavourited returns whether the video reference is favourited.
func (s *Service) VideoFavourited(ctx context.Context, videoRef models.VideoRef) (*models.Favourite, error) {
	return s.store.LookupByVideoRef(ctx, videoRef)
}

func (s *Service) FavoriteVideoByOrigin(ctx context.Context, origin models.FavouriteOrigin) (*models.Favourite, error) {
	videoRef, err := s.lookupVideoRefByOrigin(ctx, origin)
	if err != nil {
		return nil, err
	}

	f, err := s.store.LookupByVideoRef(ctx, videoRef)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get favourite for video ref: %v", videoRef)
	}
	if f != nil {
		return f, nil
	}

	return s.addFavourite(ctx, videoRef, origin)
}

func (s *Service) addFavourite(ctx context.Context, videoRef models.VideoRef, origin models.FavouriteOrigin) (*models.Favourite, error) {
	var newFavourite = &models.Favourite{
		ID:         uuid.Must(uuid.NewUUID()),
		Favourited: time.Now(),
		VideoRef:   videoRef,
	}

	switch origin.Type {
	case models.FeedItemOriginType:
		// Favourite created from feed item
		originId, err := uuid.Parse(origin.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid feed item origin ID")
		}

		feedItem, err := s.feedItemStore.Get(ctx, originId)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot get feed item with ID %v", originId)
		} else if feedItem == nil {
			return nil, errors.Wrapf(err, "feed item with ID %v is missing", originId)
		}

		newFavourite.Origin = origin
		newFavourite.Title = feedItem.Title
		newFavourite.Link = feedItem.Link
		newFavourite.Published = feedItem.Published
	}

	if err := s.store.Save(ctx, newFavourite); err != nil {
		return nil, errors.Wrapf(err, "cannot save new feed item")
	}

	return newFavourite, nil
}

func (s *Service) DeleteFavourite(ctx context.Context, favouriteId uuid.UUID) error {
	return s.store.Delete(ctx, favouriteId)
}

func (s *Service) lookupVideoRefByOrigin(ctx context.Context, origin models.FavouriteOrigin) (models.VideoRef, error) {
	switch origin.Type {
	case models.FeedItemOriginType:
		originId, err := uuid.Parse(origin.ID)
		if err != nil {
			return models.VideoRef{}, errors.Wrapf(err, "invalid feed item origin ID")
		}

		feedItem, err := s.feedItemStore.Get(ctx, originId)
		if err != nil {
			return models.VideoRef{}, errors.Wrapf(err, "cannot get feed item with ID %v", originId)
		} else if feedItem == nil {
			return models.VideoRef{}, errors.Wrapf(err, "feed item with ID %v is missing", originId)
		}

		return models.VideoRef{
			Source: models.YoutubeVideoRefSource,
			ID:     feedItem.EntryID,
		}, nil
	}

	return models.VideoRef{}, errors.Errorf("unsupported origin: %v", origin.Type)
}
