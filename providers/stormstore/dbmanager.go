package stormstore

import (
	"path/filepath"

	"github.com/asdine/storm/v3"
	"github.com/pkg/errors"
)

const (
	settingsDbFilaname = "settings.db"
)

type DBManager struct {
	baseDir   string
	databases map[string]*storm.DB
}

func NewDBManager(baseDir string) *DBManager {
	return &DBManager{
		baseDir:   baseDir,
		databases: make(map[string]*storm.DB),
	}
}

func (dm *DBManager) Close() {
	for _, d := range dm.databases {
		d.Close()
	}
}

func (dm *DBManager) Open(filename string) (*storm.DB, error) {
	db, hasDb := dm.databases[filename]
	if hasDb {
		return db, nil
	}

	fullFilename := filepath.Join(dm.baseDir, filename)
	db, err := storm.Open(fullFilename)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open storm DB: %v", fullFilename)
	}

	return db, nil
}
