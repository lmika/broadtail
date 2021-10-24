package ytdownload

import (
	"github.com/lmika/broadtail/providers/jobs"
)

type Config struct {
	LibraryDir string
}

type Service struct {
	config Config
}

func New(config Config) *Service {
	return &Service{
		config: config,
	}
}

func (s *Service) NewYoutubeDownloadTask(youtubeId string) jobs.Task {
	return &YoutubeDownloadTask{
		YoutubeId: youtubeId,
		TargetDir: s.config.LibraryDir,
	}
}
