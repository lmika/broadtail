package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/services/favourites"
	"github.com/lmika/broadtail/services/videomanager"
	"github.com/lmika/broadtail/services/videosources"
	"github.com/lmika/gopkgs/http/middleware/render"
	"github.com/pkg/errors"
)

type detailsHandler struct {
	// ytdownloadService *ytdownload.Service
	videoSources     *videosources.Service
	videoManager     *videomanager.VideoManager
	favouriteService *favourites.Service
}

func (dh *detailsHandler) QuickLook() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		youtubeId := r.FormValue("youtube_id")
		if youtubeId == "" {
			return errhandler.Errorf(http.StatusBadRequest, "missing YouTube ID")
		}

		http.Redirect(w, r, "/details/video/"+youtubeId, http.StatusSeeOther)
		return nil
	})
}

func (dh *detailsHandler) VideoDetails() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		videoID, ok := mux.Vars(r)["video_id"]
		if !ok {
			return errhandler.Errorf(http.StatusBadRequest, "invalid video ID: %v", videoID)
		}

		videoRef, err := models.ParseVideoRef(videoID)
		if err != nil {
			return errhandler.Wrap(err, http.StatusBadRequest)
		}

		videoSource, err := dh.videoSources.SourceProvider(videoRef)
		if err != nil {
			return errhandler.Wrap(err, http.StatusInternalServerError)
		}

		video, err := videoSource.GetVideoMetadata(ctx, videoRef)
		// video, err := dh.ytdownloadService.GetVideoMetadata(ctx, videoID)
		if err != nil {
			return errhandler.Wrap(err, http.StatusInternalServerError)
		}

		var downloadStatusStr string
		downloadStatus, err := dh.videoManager.DownloadStatus(videoRef)
		if err != nil {
			downloadStatusStr = "error: " + err.Error()
		} else {
			downloadStatusStr = downloadStatus.String()
		}

		favouriteStatus, err := dh.favouriteService.VideoFavourited(ctx, videoRef)
		if err != nil {
			return errors.Wrap(err, "cannot get favourite status")
		}

		if favouriteStatus != nil {
			render.Set(r, "favouriteID", favouriteStatus.ID)
		} else {
			render.Set(r, "favouriteID", "")
		}

		render.Set(r, "video", video)
		render.Set(r, "downloadStatus", downloadStatusStr)
		render.Set(r, "favouriteStatus", favouriteStatus)
		if fromFeedID := r.FormValue("from_feed_id"); fromFeedID != "" {
			render.Set(r, "fromFeedID", fromFeedID)
		}
		render.HTML(r, w, http.StatusOK, "details/show.html")
		return nil
	})
}
