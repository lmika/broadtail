package models

import (
	"os"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

type Config struct {
	MediaDir          string
	YouTubeDLCommmand string
}

func (f Config) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.MediaDir, validation.Required, validation.By(isDir)),
		validation.Field(&f.YouTubeDLCommmand, validation.Required),
	)
}

func isDir(value interface{}) error {
	str, isStr := value.(string)
	if !isStr {
		return errors.New("value is not a string")
	}

	s, err := os.Stat(str)
	if err != nil {
		return errors.Wrapf(err, "cannot stat '%v'", str)
	}

	if !s.IsDir() {
		return errors.Wrapf(err, "'%v' is not a dir")
	}
	return nil
}
