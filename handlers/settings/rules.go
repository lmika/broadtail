package settings

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/lmika/broadtail/middleware/reqbind"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/services/feedsmanager"
	"github.com/lmika/broadtail/services/rules"
	"github.com/pkg/errors"

	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/gopkgs/http/middleware/render"
)

type RulesHandlers struct {
	RulesService *rules.Service
	FeedManager  *feedsmanager.FeedsManager
}

func (sh *RulesHandlers) index() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		render.UseFrame(r, "frames/settings.html")

		rules, err := sh.RulesService.List(ctx)
		if err != nil {
			return errors.Wrap(err, "cannot get list of rules")
		}

		render.Set(r, "rules", rules)
		render.HTML(r, w, http.StatusOK, "settings/rules/index.html")
		return nil
	})
}

func (sh *RulesHandlers) newRule() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		feeds, err := sh.FeedManager.List(ctx)
		if err != nil {
			return errors.Wrap(err, "cannot list feeds")
		}

		render.Set(r, "feeds", feeds)
		render.Set(r, "rule", &models.Rule{
			Active: true,
		})
		render.Set(r, "tmplArgs", editRuleTemplateArgs{
			Path:  "/settings/rules",
			Title: "New Rule",
		})
		render.HTML(r, w, http.StatusOK, "settings/rules/edit.html")
		return nil
	})
}

func (sh *RulesHandlers) getRule() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		ruleId, err := uuid.Parse(mux.Vars(r)["id"])
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid feed ID: %v", err.Error())
		}

		rule, err := sh.RulesService.Get(ctx, ruleId)
		if err != nil {
			return errhandler.Errorf(http.StatusNotFound, "rule with ID not found: %v", ruleId)
		}

		feeds, err := sh.FeedManager.List(ctx)
		if err != nil {
			return errors.Wrap(err, "cannot list feeds")
		}

		render.Set(r, "feeds", feeds)
		render.Set(r, "rule", rule)
		render.Set(r, "tmplArgs", editRuleTemplateArgs{
			Path:  "/settings/rules/" + ruleId.String(),
			Title: "Edit Rule",
		})
		render.HTML(r, w, http.StatusOK, "settings/rules/edit.html")
		return nil
	})
}

func (sh *RulesHandlers) createRule() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var newRule models.Rule

		newRule.ID = uuid.New()
		if err := reqbind.Bind(&newRule, r); err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid request")
		}

		if err := sh.RulesService.Save(ctx, &newRule); err != nil {
			return errors.Wrap(err, "unable to save rule")
		}

		http.Redirect(w, r, "/settings/rules", http.StatusSeeOther)
		return nil
	})
}

func (sh *RulesHandlers) updateRule() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		ruleId, err := uuid.Parse(mux.Vars(r)["id"])
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid feed ID: %v", err.Error())
		}

		rule, err := sh.RulesService.Get(ctx, ruleId)
		if err != nil {
			return errhandler.Errorf(http.StatusNotFound, "rule with ID not found: %v", ruleId)
		}

		if err := reqbind.Bind(rule, r); err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid request")
		}

		if err := sh.RulesService.Save(ctx, rule); err != nil {
			return errors.Wrap(err, "unable to save rule")
		}

		http.Redirect(w, r, "/settings/rules", http.StatusSeeOther)
		return nil
	})
}

func (sh *RulesHandlers) deleteRule() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		ruleId, err := uuid.Parse(mux.Vars(r)["id"])
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid feed ID: %v", err.Error())
		}

		if err := sh.RulesService.Delete(ctx, ruleId); err != nil {
			return errors.Wrap(err, "unable to delete rule")
		}

		http.Redirect(w, r, "/settings/rules", http.StatusSeeOther)
		return nil
	})
}

func (sh *RulesHandlers) Routes(r *mux.Router) {
	r.Handle("/settings/rules", sh.index()).Methods("GET")
	r.Handle("/settings/rules", sh.createRule()).Methods("POST")
	r.Handle("/settings/rules/new", sh.newRule()).Methods("GET")
	r.Handle("/settings/rules/{id}", sh.getRule()).Methods("GET")
	r.Handle("/settings/rules/{id}", sh.updateRule()).Methods("PUT", "POST")
	r.Handle("/settings/rules/{id}", sh.deleteRule()).Methods("DELETE")
}

type editRuleTemplateArgs struct {
	Path  string
	Title string
}
