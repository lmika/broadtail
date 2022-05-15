package handlers

import (
	"context"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/services/videodownload"
	"net/http"

	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/broadtail/services/jobsmanager"
)

type youTubeDownloadHandlers struct {
	videoDownloadService *videodownload.Service
	jobsManager          *jobsmanager.JobsManager
}

func (ytdl *youTubeDownloadHandlers) CreateDownloadJob() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
		youtubeId := r.FormValue("video_ref")
		if youtubeId == "" {
			return errhandler.Errorf(http.StatusBadRequest, "missing YouTube ID")
		}

		videoRef, err := models.ParseVideoRef(youtubeId)
		if err != nil {
			return errhandler.Wrap(err, http.StatusBadRequest)
		}

		if err := ytdl.videoDownloadService.QueueForDownload(ctx, videoRef, nil); err != nil {
			return err
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	})
}
