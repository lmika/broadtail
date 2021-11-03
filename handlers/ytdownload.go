package handlers

import (
	"context"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/broadtail/middleware/render"
	"github.com/lmika/broadtail/services/jobsmanager"
	"github.com/lmika/broadtail/services/ytdownload"
	"net/http"
)

type youTubeDownloadHandlers struct {
	ytdownloadService *ytdownload.Service
	jobsManager *jobsmanager.JobsManager
}

func (ytdl *youTubeDownloadHandlers) ShowDetails() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		youtubeId := r.FormValue("youtube_id")
		if youtubeId == "" {
			return errhandler.Errorf(http.StatusBadRequest, "missing YouTube ID")
		}

		video, err := ytdl.ytdownloadService.GetVideoMetadata(ctx, youtubeId)
		if err != nil {
			return errhandler.Wrap(err, http.StatusInternalServerError)
		}

		render.Set(r, "video", video)
		if fromFeedID := r.FormValue("from_feed_id"); fromFeedID != "" {
			render.Set(r, "fromFeedID", fromFeedID)
		}
		render.HTML(r, w, http.StatusOK, "videos/show.html")
		return nil
	})
}

func (ytdl *youTubeDownloadHandlers) CreateDownloadJob() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
		youtubeId := r.FormValue("youtube_id")
		if youtubeId == "" {
			return errhandler.Errorf(http.StatusBadRequest, "missing YouTube ID")
		}

		var feedID = uuid.Nil
		if feedIDStr := r.FormValue("from_feed_id"); feedIDStr != "" {
			feedID, err = uuid.Parse(feedIDStr)
			if err != nil {
				return err
			}
		}

		task, err := ytdl.ytdownloadService.NewYoutubeDownloadTask(ctx, youtubeId, feedID)
		if err != nil {
			return err
		}

		ytdl.jobsManager.Dispatcher().Enqueue(task)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	})
}
