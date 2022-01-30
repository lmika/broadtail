package videomanager

import "github.com/lmika/broadtail/models"

type VideoStore interface {
	ListRecent() ([]models.SavedVideo, error)
	FindWithExtID(id string) (*models.SavedVideo, error)
}
