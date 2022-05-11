package youtuberss

import (
	"context"
	"encoding/xml"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/models/ytrss"
	"github.com/pkg/errors"
)

type Provider struct {
	client *resty.Client
}

func New() *Provider {
	return &Provider{client: resty.New()}
}

func (p *Provider) GetForFeed(ctx context.Context, feed models.Feed) ([]models.FetchedFeedItem, error) {
	ytrssEntries, err := p.getForFeed(ctx, feed)
	if err != nil {
		return nil, err
	}
	return p.convertToFetchedFeedItem(ytrssEntries), nil
}

func (fm *Provider) FeedExternalURL(f models.Feed) (string, error) {
	switch f.Type {
	case models.FeedTypeYoutubeChannel:
		return fmt.Sprintf("https://www.youtube.com/channel/%v", f.ExtID), nil
	case models.FeedTypeYoutubePlaylist:
		return fmt.Sprintf("https://www.youtube.com/playlist/%v", f.ExtID), nil
	}
	return "", errors.Errorf("external url unsupported for feed type: %v", f.Type)
}

func (fm *Provider) FeedHints(feed models.Feed) models.FeedHints {
	return models.FeedHints{
		Ordering: models.ChronologicalFeedItemOrdering,
	}
}

func (p *Provider) getForFeed(ctx context.Context, feed models.Feed) ([]ytrss.Entry, error) {
	switch feed.Type {
	case models.FeedTypeYoutubeChannel:
		return p.getForChannelID(ctx, feed.ExtID)
	case models.FeedTypeYoutubePlaylist:
		return p.getForPlaylistID(ctx, feed.ExtID)
	}
	return nil, errors.Errorf("unrecognised feed type: %v", feed.Type)
}

func (p *Provider) convertToFetchedFeedItem(rssEntries []ytrss.Entry) []models.FetchedFeedItem {
	ffis := make([]models.FetchedFeedItem, len(rssEntries))
	for i, rssEntry := range rssEntries {
		ffis[i] = models.FetchedFeedItem{
			VideoRef: models.VideoRef{
				Source: models.YoutubeVideoRefSource,
				ID:     rssEntry.VideoID,
			},
			Title:     rssEntry.Title,
			Link:      rssEntry.Link,
			Published: rssEntry.Published,
		}
	}
	return ffis
}

func (p *Provider) getForChannelID(ctx context.Context, channelID string) ([]ytrss.Entry, error) {
	return p.getFeed(p.client.R().
		SetQueryParam("channel_id", channelID).
		Get("https://www.youtube.com/feeds/videos.xml"))
}

func (p *Provider) getForPlaylistID(ctx context.Context, playlistID string) ([]ytrss.Entry, error) {
	return p.getFeed(p.client.R().
		SetQueryParam("playlist_id", playlistID).
		Get("https://www.youtube.com/feeds/videos.xml"))
}

func (p *Provider) getFeed(resp *resty.Response, err error) ([]ytrss.Entry, error) {
	if err != nil {
		return nil, errors.Wrapf(err, "error getting RSS feed")
	} else if !resp.IsSuccess() {
		return nil, errors.Errorf("error getting RSS feed: HTTP code %v", resp.StatusCode())
	}

	var feed ytrss.Feed
	if err := xml.Unmarshal(resp.Body(), &feed); err != nil {
		return nil, errors.Wrapf(err, "cannot marshal XML response body")
	}

	return feed.Entries, nil
}
