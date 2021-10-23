package handlers

import (
	"context"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/broadtail/middleware/jobdispatcher"
	"github.com/lmika/broadtail/services/ytdownload"
	"net/http"
)

type youTubeDownloadHandlers struct {
	ytdownloadService *ytdownload.Service
}

func (ytdl *youTubeDownloadHandlers) CreateDownloadJob() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		youtubeId := r.FormValue("youtube_id")
		if youtubeId == "" {
			return errhandler.Errorf(http.StatusBadRequest, "missing YouTube ID")
		}


		task := ytdl.ytdownloadService.NewYoutubeDownloadTask(youtubeId)
		jobdispatcher.FromContext(r.Context()).Dispatcher.Enqueue(task)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	})
}
