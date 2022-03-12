package stormstore

import (
	"context"
	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/pkg/errors"
)

type FavouriteStore struct {
	db *storm.DB
}

func NewFavouriteStore(filename string) (*FavouriteStore, error) {
	db, err := storm.Open(filename)
	if err != nil {
		return nil, err
	}

	return &FavouriteStore{db: db}, nil
}

func (f *FavouriteStore) Close() {
	f.db.Close()
}

func (fs *FavouriteStore) LookupByVideoRef(ctx context.Context, videoRef models.VideoRef) (*models.Favourite, error) {
	var favourite models.Favourite
	if err := fs.db.One("VideoRef", videoRef, &favourite); err != nil {
		if err == storm.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &favourite, nil
}

func (f *FavouriteStore) Save(ctx context.Context, favourite *models.Favourite) error {
	return f.db.Save(favourite)
}

func (f *FavouriteStore) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.Wrapf(f.db.DeleteStruct(&models.Favourite{ID: id}), "cannot delete feed item: %v", id)
}
