package rules

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/pkg/errors"
)

type Service struct {
	store     RuleStore
	feedStore FeedStore
}

func NewService(store RuleStore, feedStore FeedStore) *Service {
	return &Service{
		store:     store,
		feedStore: feedStore,
	}
}

func (s *Service) List(ctx context.Context) ([]*models.RuleWithDescription, error) {
	rules, err := s.store.List(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot list rules")
	}

	return mapSlice(rules, func(rule *models.Rule) *models.RuleWithDescription {
		return &models.RuleWithDescription{Rule: *rule, Description: s.describeRule(ctx, rule)}
	}), nil
}

func (s *Service) Get(ctx context.Context, ruleID uuid.UUID) (*models.Rule, error) {
	return s.store.Get(ctx, ruleID)
}

func (s *Service) Save(ctx context.Context, rule *models.Rule) error {
	return s.store.Save(ctx, rule)
}

func (s *Service) Delete(ctx context.Context, ruleID uuid.UUID) error {
	return s.store.Delete(ctx, ruleID)
}

func (s *Service) describeRule(ctx context.Context, rule *models.Rule) string {
	// Collect conditions
	conditionDescription := make([]string, 0)
	if rule.Condition.FeedID != uuid.Nil {
		if feed, err := s.feedStore.Get(ctx, rule.Condition.FeedID); err == nil {
			conditionDescription = append(conditionDescription, fmt.Sprintf("feed is '%v'", feed.Name))
		} else {
			conditionDescription = append(conditionDescription, "feed is '(unknown)'")
		}
	}
	if rule.Condition.Title != "" {
		conditionDescription = append(conditionDescription, fmt.Sprintf("title is '%v'", rule.Condition.Title))
	}
	if rule.Condition.Description != "" {
		conditionDescription = append(conditionDescription, fmt.Sprintf("description is '%v'", rule.Condition.Description))
	}

	// Collect actions
	actionDescription := make([]string, 0)
	if rule.Action.Download {
		actionDescription = append(actionDescription, "download the video")
	}
	if rule.Action.MarkFavourite {
		actionDescription = append(actionDescription, "mark as favourite")
	}

	fullDescription := strings.Builder{}
	if len(conditionDescription) > 0 {
		fullDescription.WriteString("WHEN ")
		fullDescription.WriteString(strings.Join(conditionDescription, " and "))
		fullDescription.WriteString(",\nTHEN ")
	} else {
		fullDescription.WriteString("For every feed item, ")
	}
	if len(actionDescription) > 0 {
		fullDescription.WriteString(strings.Join(actionDescription, " and "))
	} else {
		fullDescription.WriteString("do nothing")
	}
	fullDescription.WriteString(".")

	return fullDescription.String()
}

func mapSlice[T, U any](ts []T, mapper func(t T) U) []U {
	us := make([]U, len(ts))
	for i, t := range ts {
		us[i] = mapper(t)
	}
	return us
}
