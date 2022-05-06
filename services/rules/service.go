package rules

import (
	"context"

	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/pkg/errors"
)

type Service struct {
	store RuleStore
}

func NewService(store RuleStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) List(ctx context.Context) ([]*models.RuleWithDescription, error) {
	rules, err := s.store.List(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot list rules")
	}

	return mapSlice(rules, func(rule *models.Rule) *models.RuleWithDescription {
		return &models.RuleWithDescription{Rule: *rule, Description: s.describeRule(rule)}
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

func (s *Service) describeRule(rule *models.Rule) string {
	return "TODO"
}

func mapSlice[T, U any](ts []T, mapper func(t T) U) []U {
	us := make([]U, len(ts))
	for i, t := range ts {
		us[i] = mapper(t)
	}
	return us
}
