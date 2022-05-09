package videosources

import "github.com/lmika/broadtail/models"

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

// SourceProvider returns the SourceProvider that is registered to handle videos from the given source.
func (s *Service) SourceProvider(videoRef models.VideoRef) (SourceProvider, error) {
	return nil, nil
}
