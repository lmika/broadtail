package videosources

import (
	"context"

	"github.com/lmika/broadtail/models"
)

// SourceProvider is a provider for a particular video source.
type SourceProvider interface {
	// GetVideoMetadata returns the metadata details for a sourceID.
	GetVideoMetadata(ctx context.Context, videoRef models.VideoRef) (*models.Video, error)

	GetVideoURL(videoRef models.VideoRef) string

	// DownloadVideo will download the video.
	DownloadVideo(ctx context.Context, videoRef models.VideoRef, options models.DownloadOptions, logline func(line models.LogMessage)) (outputFilename string, err error)
}
