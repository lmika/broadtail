package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/broadtail/middleware/render"
	"github.com/lmika/broadtail/services/ytdownload"
	"net/http"
)

type detailsHandler struct {
	ytdownloadService *ytdownload.Service
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

		video, err := dh.ytdownloadService.GetVideoMetadata(ctx, videoID)
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
