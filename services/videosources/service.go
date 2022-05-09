package videosources

import "github.com/lmika/broadtail/models"

type Service struct {
	provider SourceProvider
}

func NewService(provider SourceProvider) *Service {
	return &Service{
		provider: provider,
	}
}

// SourceProvider returns the SourceProvider that is registered to handle videos from the given source.
func (s *Service) SourceProvider(videoRef models.VideoRef) (SourceProvider, error) {
	return nil, nil
}
