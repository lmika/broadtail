package monitor

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/broadtail/services/videomanager"
	"github.com/lmika/gopkgs/http/middleware/render"
)

type VideoHandlers struct {
	VideoManager *videomanager.VideoManager
}

func (h *VideoHandlers) List() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		render.UseFrame(r, "frames/history.html")

		videoList, err := h.VideoManager.List()
		if err != nil {
			return err
		}

		render.Set(r, "videos", videoList)
		render.HTML(r, w, http.StatusOK, "videos/index.html")
		return nil
	})
}

func (h *VideoHandlers) Show() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		feedId, err := uuid.Parse(mux.Vars(r)["video_id"])
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid feed ID: %v", err.Error())
		}

		savedVideo, err := h.VideoManager.Get(feedId)
		if err != nil {
			return errhandler.Errorf(http.StatusInternalServerError, "cannot get saved video - %v: %v", feedId.String(), err)
		}

		render.Set(r, "video", savedVideo)
		render.HTML(r, w, http.StatusOK, "videos/show.html")
		return nil
	})
}
