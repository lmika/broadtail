package rssfetcher

import (
	"context"
	"encoding/xml"
	"github.com/go-resty/resty/v2"
	"github.com/lmika/broadtail/models/ytrss"
	"github.com/pkg/errors"
)

type Provider struct {
	client *resty.Client
}

func New() *Provider {
	return &Provider{client: resty.New()}
}

func (p *Provider) GetForChannelID(ctx context.Context, channelID string) ([]ytrss.Entry, error) {
	return p.getFeed(p.client.R().
		SetQueryParam("channel_id", channelID).
		Get("https://www.youtube.com/feeds/videos.xml"))
}

func (p *Provider) GetForPlaylistID(ctx context.Context, playlistID string) ([]ytrss.Entry, error) {
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