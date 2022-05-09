package favourites

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/pkg/errors"
)

type Service struct {
	store         FavouriteStore
	videoMetadata VideoMetadata
	feedItemStore FeedItemStore
	feedStore     FeedStore
}

func NewService(store FavouriteStore, videoMetadata VideoMetadata, feedStore FeedStore, feedItemStore FeedItemStore) *Service {
	return &Service{
		store:         store,
		videoMetadata: videoMetadata,
		feedStore:     feedStore,
		feedItemStore: feedItemStore,
	}
}

func (s *Service) List(ctx context.Context, query models.FeedItemFilter, page int) ([]models.FavouriteWithOrigin, error) {
	favourites, err := s.store.List(ctx, query, page)
	if err != nil {
		return nil, err
	}

	var seenFeeds = make(map[uuid.UUID]models.Feed)

	favouritesWithOrigin := make([]models.FavouriteWithOrigin, len(favourites))
	for i, f := range favourites {
		var favouriteWithOrigin = models.FavouriteWithOrigin{Favourite: f}

		switch f.Origin.Type {
		case models.FeedItemOriginType:
			if feedItem, err := s.feedItemStore.Get(ctx, uuid.MustParse(f.Origin.ID)); err == nil {
				feed, ok := seenFeeds[feedItem.FeedID]
				if !ok {
					if feed, err = s.feedStore.Get(ctx, feedItem.FeedID); err == nil {
						seenFeeds[feedItem.FeedID] = feed
					} else {
						return nil, errors.Wrapf(err, "cannot get feed with ID %v", feedItem.FeedID)
					}
				}

				if feed.ID != uuid.Nil {
					favouriteWithOrigin.OriginTitle = feed.Name
					favouriteWithOrigin.OriginURL = "/feeds/" + feed.ID.String()
				} else {
					favouriteWithOrigin.OriginTitle = "(unknown)"
				}
			} else {
				favouriteWithOrigin.OriginTitle = "(unknown)"
			}
		case models.ManualOriginType:
			favouriteWithOrigin.OriginTitle = "Manual"
		default:
			favouriteWithOrigin.OriginTitle = "(unknown)"
		}

		favouritesWithOrigin[i] = favouriteWithOrigin
	}

	return favouritesWithOrigin, nil
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
	case models.ManualOriginType:
		videoRef, err := models.ParseVideoRef(origin.ID)
		if err != nil {
			return nil, err
		}

		// TODO: this only supports YouTube
		videoMetadata, err := s.videoMetadata.GetVideoMetadata(ctx, videoRef.ID)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot get video metadata %v", videoRef.ID)
		}

		newFavourite.Origin = origin
		newFavourite.Title = videoMetadata.Title
		newFavourite.Link = fmt.Sprintf("https://www.youtube.com/watch?v=%v", videoMetadata.ExtID)
		newFavourite.Published = videoMetadata.UploadedOn
		newFavourite.VideoRef = videoRef
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

		return feedItem.VideoRef(), nil
		// return models.VideoRef{
		// 	Source: models.YoutubeVideoRefSource,
		// 	ID:     feedItem.EntryID,
		// }, nil
	case models.ManualOriginType:
		return models.VideoRef{Source: models.YoutubeVideoRefSource, ID: origin.ID}, nil

	}

	return models.VideoRef{}, errors.Errorf("unsupported origin: %v", origin.Type)
}
