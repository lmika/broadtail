package feedfetchers

import (
	"context"

	"github.com/lmika/broadtail/models"
	"github.com/pkg/errors"
)

type Service struct {
	drivers map[string]FeedDriver
}

func NewService(drivers map[string]FeedDriver) *Service {
	return &Service{
		drivers: drivers,
	}
}

func (s *Service) GetForFeed(ctx context.Context, feed models.Feed) ([]models.FetchedFeedItem, error) {
	feedDriver, hasDriver := s.drivers[feed.Type]
	if !hasDriver {
		return nil, errors.Errorf("missing driver for feed-type: %v", feed.Type)
	}

	return feedDriver.GetForFeed(ctx, feed)
}

func (s *Service) FeedExternalURL(feed models.Feed) (string, error) {
	feedDriver, hasDriver := s.drivers[feed.Type]
	if !hasDriver {
		return "", errors.Errorf("missing driver for feed-type: %v", feed.Type)
	}

	return feedDriver.FeedExternalURL(feed)
}

func (s *Service) FeedHints(feed models.Feed) models.FeedHints {
	if feedDriver, hasDriver := s.drivers[feed.Type]; hasDriver {
		return feedDriver.FeedHints(feed)
	}

	return models.FeedHints{
		Ordering: models.ChronologicalFeedItemOrdering,
	}
}
