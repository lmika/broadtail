package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/gopkgs/http/middleware/render"
)

type settingHandlers struct{}

func (settingHandlers) index() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		render.UseFrame(r, "frames/settings.html")

		render.HTML(r, w, http.StatusOK, "settings/general.html")
		return nil
	})
}

func (settingHandlers) rules() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		render.UseFrame(r, "frames/settings.html")

		render.HTML(r, w, http.StatusOK, "settings/rules/index.html")
		return nil
	})
}

func (settingHandlers) newRule() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		render.HTML(r, w, http.StatusOK, "settings/rules/new.html")
		return nil
	})
}

func (sh *settingHandlers) Routes(r *mux.Router) {
	r.Handle("/settings", sh.index()).Methods("GET")
	r.Handle("/settings/rules", sh.rules()).Methods("GET")
	r.Handle("/settings/rules/new", sh.newRule()).Methods("GET")
}
