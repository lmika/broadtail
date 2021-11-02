package reqbind

import (
	"github.com/lmika/broadtail/middleware/errhandler"
	"net/http"
)

func doValidate(target interface{}, r *http.Request) error {
	validatable, isValidatable := target.(Validatable)
	if !isValidatable {
		return nil
	}

	return errhandler.Wrap(validatable.Validate(), http.StatusBadRequest)
}

type Validatable interface {
	Validate() error
}

