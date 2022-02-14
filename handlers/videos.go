package handlers

import (
	"context"
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
