package youtubevideosource

import (
	"context"
	"github.com/lmika/broadtail/models"
)

type DownloadProvider interface {
	GetVideoMetadata(ctx context.Context, youtubeId string) (*models.Video, error)
	DownloadVideo(ctx context.Context, youtubeId string, options models.DownloadOptions, logline func(line models.LogMessage)) (outputFilename string, err error)
}
