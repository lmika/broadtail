package youtubevideosource

import (
	"context"
	"github.com/lmika/broadtail/models"
)

type Service struct {
	provider DownloadProvider
}

func NewService(provider DownloadProvider) *Service {
	return &Service{
		provider: provider,
	}
}

func (s *Service) GetVideoMetadata(ctx context.Context, videoRef models.VideoRef) (*models.Video, error) {
	return s.provider.GetVideoMetadata(ctx, videoRef.ID)
}

func (s *Service) DownloadVideo(ctx context.Context, videoRef models.VideoRef, options models.DownloadOptions, logline func(line string)) (outputFilename string, err error) {
	return s.provider.DownloadVideo(ctx, videoRef.ID, options, logline)
}
