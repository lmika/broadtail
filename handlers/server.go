package handlers

import (
	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/jobs"
	"github.com/lmika/broadtail/middleware/jobdispatcher"
	"github.com/lmika/broadtail/middleware/render"
	"github.com/pkg/errors"
	"html/template"
	"net/http"
	"os"
)

func Server() (http.Handler, error) {
	dispatcher := jobs.New()
	tmpls, err := template.ParseFS(os.DirFS("templates"), "*.html")
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse templates")
	}

	r := mux.NewRouter()
	r.Handle("/", indexHandler()).Methods("GET")

	handler := jobdispatcher.New(dispatcher).Use(r)
	handler = render.New(tmpls).Use(handler)

	return handler, nil
}
