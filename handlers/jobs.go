package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/broadtail/middleware/jobdispatcher"
	"github.com/lmika/broadtail/services/jobsmanager"
	"github.com/lmika/gopkgs/http/middleware/render"
)

type jobsHandlers struct {
	jobsManager *jobsmanager.JobsManager
}

func (h *jobsHandlers) Delete() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		jobIdStr := mux.Vars(r)["job_id"]
		jobId, err := uuid.Parse(jobIdStr)
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid job ID: %v", err.Error())
		}

		job := h.jobsManager.Dispatcher().Job(jobId)
		if job == nil {
			return errhandler.Errorf(http.StatusNotFound, "job %v not found", jobId)
		}
		job.Cancel()

		// Flash?
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	})
}

func (ytdl *jobsHandlers) ClearDone() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		dispatcher := jobdispatcher.FromContext(ctx).Dispatcher
		dispatcher.ClearDone()

		// Flash?
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	})
}

func (h *jobsHandlers) List() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		historialJobs, err := h.jobsManager.HistoricalJobs()
		if err != nil {
			return err
		}

		render.Set(r, "runningJobs", h.jobsManager.RecentJobs())
		render.Set(r, "historicalJobs", historialJobs)
		render.HTML(r, w, http.StatusOK, "jobs/index.html")
		return nil
	})
}

func (h *jobsHandlers) Show() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		jobIdStr := mux.Vars(r)["job_id"]
		jobId, err := uuid.Parse(jobIdStr)
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid job ID: %v", err.Error())
		}

		job, err := h.jobsManager.Job(jobId)
		if err != nil {
			return err
		}

		render.Set(r, "job", job)
		render.HTML(r, w, http.StatusOK, "jobs/show.html")
		return nil
	})
}
