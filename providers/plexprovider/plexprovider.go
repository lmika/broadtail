package plexprovider

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type PlexProvider struct {
	client    *resty.Client
	plexToken string
}

func New(basePath, plexToken string) *PlexProvider {
	var client *resty.Client = nil

	if basePath != "" {
		client = resty.New().SetBaseURL(basePath)
	}

	return &PlexProvider{
		client:    client,
		plexToken: plexToken,
	}
}

func (p *PlexProvider) NewVideoDownloaded(ctx context.Context, outVideoFile string) error {
	if p.client == nil {
		return nil
	}

	resp, err := p.client.R().
		SetQueryParam("X-Plex-Token", p.plexToken).
		Get("/library/sections/all/refresh")
	if err != nil {
		return errors.Wrap(err, "cannot refresh library")
	} else if resp.IsSuccess() {
		return errors.Errorf("cannot refresh library: non-200 response: %d", resp.StatusCode())
	}
	return nil
}
