package handlers

import (
	"context"
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
		render.HTML(r, w, http.StatusOK, "videos/show.html")
		return nil
	})
}

func (ytdl *youTubeDownloadHandlers) CreateDownloadJob() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		youtubeId := r.FormValue("youtube_id")
		if youtubeId == "" {
			return errhandler.Errorf(http.StatusBadRequest, "missing YouTube ID")
		}

		task := ytdl.ytdownloadService.NewYoutubeDownloadTask(youtubeId)
		ytdl.jobsManager.Dispatcher().Enqueue(task)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	})
}
