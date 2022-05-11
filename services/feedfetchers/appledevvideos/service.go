package appledevvideos

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/lmika/broadtail/models"
	"github.com/pkg/errors"
	"log"
	"strings"
)

type Service struct {
	client *resty.Client
}

func NewService() *Service {
	return &Service{
		client: resty.New(),
	}
}

func (s *Service) GetForFeed(ctx context.Context, feed models.Feed) ([]models.FetchedFeedItem, error) {
	resp, err := s.client.R().SetPathParam("extId", feed.ExtID).Get("https://developer.apple.com/videos/{extId}")
	if err != nil {
		return nil, errors.Wrap(err, "cannot make request to developer.apple.com")
	} else if !resp.IsSuccess() {
		return nil, errors.Wrapf(err, "non-200 response code from developer.apple.com: %v", resp.Status())
	}

	dom, err := goquery.NewDocumentFromReader(bytes.NewBuffer(resp.Body()))
	if err != nil {
		return nil, errors.Wrap(err, "cannot read dom from page")
	}

	var fetchedFeedItem = make([]models.FetchedFeedItem, 0)
	dom.Find("ul.collection-items a[href^='/videos/play']").Each(func(i int, selection *goquery.Selection) {
		h4s := selection.Find("h4")
		if h4s.Length() != 1 {
			return
		}

		refId := strings.TrimSuffix(strings.TrimPrefix(selection.AttrOr("href", ""), "/videos/play/"), "/")
		refId = strings.ReplaceAll(refId, "/", ".")
		if refId == "" {
			return
		}

		log.Println(refId)
		fetchedFeedItem = append(fetchedFeedItem, models.FetchedFeedItem{
			VideoRef: models.VideoRef{
				Source: models.AppleDevVideoRefSource,
				ID:     refId,
			},
			Title: strings.TrimSpace(h4s.First().Text()),
		})
	})

	return fetchedFeedItem, nil
}

func (s *Service) FeedExternalURL(feed models.Feed) (string, error) {
	return fmt.Sprintf("https://developer.apple.com/videos/%v", feed.ExtID), nil
}

func (fm *Service) FeedHints(feed models.Feed) models.FeedHints {
	return models.FeedHints{
		Ordering: models.AlphabeticalFeedItemOrdering,
	}
}
