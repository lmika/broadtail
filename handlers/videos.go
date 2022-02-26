package handlers

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/broadtail/middleware/render"
	"github.com/lmika/broadtail/services/videomanager"
	"net/http"
)

type videoHandlers struct {
	videoManager *videomanager.VideoManager
}

func (h *videoHandlers) List() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		videoList, err := h.videoManager.List()
		if err != nil {
			return err
		}

		render.Set(r, "videos", videoList)
		render.HTML(r, w, http.StatusOK, "videos/index.html")
		return nil
	})
}

func (h *videoHandlers) Show() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		feedId, err := uuid.Parse(mux.Vars(r)["video_id"])
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid feed ID: %v", err.Error())
		}

		savedVideo, err := h.videoManager.Get(feedId)
		if err != nil {
			return errhandler.Errorf(http.StatusInternalServerError, "cannot get saved video - %v: %v", feedId.String(), err)
		}

		render.Set(r, "video", savedVideo)
		render.HTML(r, w, http.StatusOK, "videos/show.html")
		return nil
	})
}
