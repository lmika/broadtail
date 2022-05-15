package videosources

import (
	"github.com/lmika/broadtail/models"
	"github.com/pkg/errors"
)

type Service struct {
	providers map[models.VideoRefSource]SourceProvider
}

func NewService(providers map[models.VideoRefSource]SourceProvider) *Service {
	return &Service{
		providers: providers,
	}
}

// SourceProvider returns the SourceProvider that is registered to handle videos from the given source.
func (s *Service) SourceProvider(videoRef models.VideoRef) (SourceProvider, error) {
	p, hasP := s.providers[videoRef.Source]
	if !hasP {
		return nil, errors.Errorf("unrecognised video source: %v", videoRef.Source)
	}
	return p, nil
}
