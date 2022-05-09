package videosources

import (
	"context"

	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/providers/jobs"
)

// SourceProvider is a provider for a particular video source.
type SourceProvider interface {
	// GetVideoMetadata returns the metadata details for a sourceID.
	GetVideoMetadata(ctx context.Context, videoRef models.VideoRef) (*models.Video, error)

	// BuildDownloadTask returns a new task that will download the video.
	BuildDownloadTask(ctx context.Context, videoRef models.VideoRef, options models.DownloadOptions) (jobs.Task, error)
}
