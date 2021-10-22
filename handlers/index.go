package handlers

import (
	"github.com/lmika/broadtail/middleware/jobdispatcher"
	"github.com/lmika/broadtail/middleware/render"
	"net/http"
)

func indexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobs := jobdispatcher.FromContext(r.Context()).Dispatcher.List()

		render.Set(r, "jobs", jobs)
		render.HTML(r, w, http.StatusOK, "index.html")
	})
}
