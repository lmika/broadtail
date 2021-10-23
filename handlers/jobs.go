package handlers

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/broadtail/middleware/jobdispatcher"
	"net/http"
)

type jobsHandlers struct {
}

func (ytdl *jobsHandlers) Delete() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		jobIdStr := mux.Vars(r)["job_id"]
		jobId, err := uuid.Parse(jobIdStr)
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid job ID: %v", err.Error())
		}

		dispatcher := jobdispatcher.FromContext(ctx).Dispatcher

		job := dispatcher.Job(jobId)
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