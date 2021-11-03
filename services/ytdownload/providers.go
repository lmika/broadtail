package ytdownload

import (
	"context"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
)

type Provider interface {
	GetVideoMetadata(ctx context.Context, youtubeId string) (*models.Video, error)
	DownloadVideo(ctx context.Context, options models.DownloadOptions, logline func(line string)) error
}

type FeedStore interface {
	Get(ctx context.Context, id uuid.UUID) (models.Feed, error)
}