package handlers

import (
	"github.com/lmika/broadtail/middleware/jobdispatcher"
	"github.com/lmika/broadtail/services/ytdownload"
	"net/http"
)

type youTubeDownloadHandlers struct {
	ytdownloadService *ytdownload.Service
}

func (ytdl *youTubeDownloadHandlers) CreateDownloadJob() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task := ytdl.ytdownloadService.NewYoutubeDownloadTask("SSS")
		jobdispatcher.FromContext(r.Context()).Dispatcher.Enqueue(task)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
