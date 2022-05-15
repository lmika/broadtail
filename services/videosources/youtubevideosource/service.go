package youtubevideosource

import (
	"context"
	"fmt"
	"github.com/lmika/broadtail/models"
	"net/url"
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

func (s *Service) DownloadVideo(ctx context.Context, videoRef models.VideoRef, options models.DownloadOptions, logline func(line models.LogMessage)) (outputFilename string, err error) {
	return s.provider.DownloadVideo(ctx, videoRef.ID, options, logline)
}

func (s *Service) GetVideoURL(videoRef models.VideoRef) string {
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", url.QueryEscape(videoRef.ID))
}
