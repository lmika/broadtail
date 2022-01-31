package handlers

import (
	"context"
	"net/http"

	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/broadtail/middleware/render"
	"github.com/lmika/broadtail/middleware/reqbind"
	"github.com/lmika/broadtail/middleware/sessions"
	"github.com/lmika/broadtail/models"
)

type configHandler struct{}

func (cf *configHandler) Show() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		config := models.Config{
			MediaDir:          "/tmp",
			YouTubeDLCommmand: "youtube-dl",
		}

		render.Set(r, "config", config)
		render.HTML(r, w, http.StatusOK, "config/index.html")
		return nil
	})
}

func (cf *configHandler) Update() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var config models.Config

		if err := reqbind.Bind(&config, r); err != nil {
			return err
		}

		// TEMP
		sess, err := sessions.Store(r).Get(r, "broadtail")
		if err != nil {
			return err
		}

		render.Set(r, "config", config)
		render.HTML(r, w, http.StatusOK, "config/index.html")
		return nil
	})
}
