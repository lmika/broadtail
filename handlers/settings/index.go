package settings

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/gopkgs/http/middleware/render"
	"net/http"
)

type IndexHandlers struct {
}

func (IndexHandlers) index() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		render.UseFrame(r, "frames/settings.html")

		render.HTML(r, w, http.StatusOK, "settings/general.html")
		return nil
	})
}

func (sh *IndexHandlers) Routes(r *mux.Router) {
	r.Handle("/settings", sh.index()).Methods("GET")
}
