package handlers

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/lmika/broadtail/middleware/reqbind"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/services/feedsmanager"
	"github.com/pkg/errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/gopkgs/http/middleware/render"
)

type settingHandlers struct {
	feedManager *feedsmanager.FeedsManager
}

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

func (sh *settingHandlers) newRule() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		feeds, err := sh.feedManager.List(ctx)
		if err != nil {
			return errors.Wrap(err, "cannot list feeds")
		}

		render.Set(r, "feeds", feeds)
		render.HTML(r, w, http.StatusOK, "settings/rules/new.html")
		return nil
	})
}

func (settingHandlers) createRule() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var newRule models.Rule

		if err := reqbind.Bind(&newRule, r); err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid request")
		}
		spew.Dump(newRule)

		http.Redirect(w, r, "/settings/rules", http.StatusSeeOther)
		return nil
	})
}

func (sh *settingHandlers) Routes(r *mux.Router) {
	r.Handle("/settings", sh.index()).Methods("GET")
	r.Handle("/settings/rules", sh.rules()).Methods("GET")
	r.Handle("/settings/rules", sh.createRule()).Methods("POST")
	r.Handle("/settings/rules/new", sh.newRule()).Methods("GET")
}
