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

//func (ytdl *youTubeDownloadHandlers) ShowDetails() http.Handler {
//	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
//		youtubeId := r.FormValue("video_ref")
//		if youtubeId == "" {
//			return errhandler.Errorf(http.StatusBadRequest, "missing YouTube ID")
//		}
//
//		videoRef, err := models.ParseVideoRef(youtubeId)
//		if err != nil {
//			return errhandler.Wrap(err, http.StatusBadRequest)
//		}
//
//		video, err := ytdl.videoDownloadService.GetVideoMetadata(ctx, videoRef)
//		if err != nil {
//			return errhandler.Wrap(err, http.StatusInternalServerError)
//		}
//
//		render.Set(r, "video", video)
//		render.HTML(r, w, http.StatusOK, "videos/show.html")
//		return nil
//	})
//}

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

		//var feedID = uuid.Nil
		//if feedIDStr := r.FormValue("from_feed_id"); feedIDStr != "" {
		//	feedID, err = uuid.Parse(feedIDStr)
		//	if err != nil {
		//		return err
		//	}
		//}

		if err := ytdl.videoDownloadService.QueueForDownload(ctx, videoRef, nil); err != nil {
			return err
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	})
}
