package videodownload

import (
	"context"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
)

type DownloadProvider interface {
	GetVideoMetadata(ctx context.Context, videoRef models.VideoRef) (*models.Video, error)
	DownloadVideo(ctx context.Context, videoRef models.VideoRef, options models.DownloadOptions, logline func(line string)) (outputFilename string, err error)
}

type FeedStore interface {
	Get(ctx context.Context, id uuid.UUID) (models.Feed, error)
}

type FeedItemStore interface {
	GetByVideoRef(ctx context.Context, videoRef models.VideoRef) (*models.FeedItem, error)
}

type VideoStore interface {
	Save(video models.SavedVideo) error
	DeleteWithExtID(extId models.VideoRef) error
}

type VideoDownloadHooks interface {
	NewVideoDownloaded(ctx context.Context, outVideoFile string) error
}
