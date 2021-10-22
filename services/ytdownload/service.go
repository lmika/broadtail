package ytdownload

import "github.com/lmika/broadtail/jobs"

type Service struct {

}

func New() *Service {
	return &Service{}
}

func (s *Service) NewYoutubeDownloadTask(youtubeId string) jobs.Task {
	return &YoutubeDownloadTask{}
}
