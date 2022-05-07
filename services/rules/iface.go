package rules

import (
	"context"

	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
)

type RuleStore interface {
	List(ctx context.Context) ([]*models.Rule, error)
	Save(ctx context.Context, rule *models.Rule) error
	Get(ctx context.Context, ruleID uuid.UUID) (*models.Rule, error)
	Delete(ctx context.Context, ruleID uuid.UUID) error
}

type FeedStore interface {
	Get(ctx context.Context, id uuid.UUID) (models.Feed, error)
}
