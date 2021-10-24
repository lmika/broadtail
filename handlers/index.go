package handlers

import (
	"github.com/lmika/broadtail/middleware/render"
	"github.com/lmika/broadtail/services/jobsmanager"
	"net/http"
)

type indexHandlers struct {
	jobsManager *jobsmanager.JobsManager
}

func (ih *indexHandlers) Index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.Set(r, "jobs", ih.jobsManager.RecentJobs())
		render.HTML(r, w, http.StatusOK, "index.html")
	})
}
