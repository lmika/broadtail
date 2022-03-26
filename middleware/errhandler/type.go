package errhandler

import (
	"context"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

type ErrorHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

func HandlerFunc(errHandler ErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := errHandler(r.Context(), w, r)
		if err != nil {
			var status = http.StatusInternalServerError
			var herr handlerError
			if errors.Is(err, &herr) {
				status = herr.httpStatus
			}

			log.Printf("error: %v", err)
			http.Error(w, err.Error(), status)
		}
	})
}

func Errorf(status int, fmt string, args ...interface{}) error {
	return handlerError{httpStatus: status, cause: errors.Errorf(fmt, args...)}
}

func Wrap(cause error, status int) error {
	if cause == nil {
		return cause
	}
	return handlerError{httpStatus: status, cause: cause}
}

type handlerError struct {
	httpStatus int
	cause      error
}

func (eh handlerError) Error() string {
	return eh.cause.Error()
}

func (eh handlerError) Unwrap() error {
	return eh.cause
}
