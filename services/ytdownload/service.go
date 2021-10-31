package ytdownload

import (
	"context"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/providers/jobs"
)

type Config struct {
	LibraryDir string
}

type Service struct {
	config Config
	provider Provider
}

func New(config Config, provider Provider) *Service {
	return &Service{
		config: config,
		provider: provider,
	}
}

func (s *Service) GetVideoMetadata(ctx context.Context, youtubeId string) (*models.Video, error) {
	return s.provider.GetVideoMetadata(ctx, youtubeId)
}

func (s *Service) NewYoutubeDownloadTask(youtubeId string) jobs.Task {
	return &YoutubeDownloadTask{
		Provider: s.provider,
		YoutubeId: youtubeId,
		TargetDir: s.config.LibraryDir,
	}
}
