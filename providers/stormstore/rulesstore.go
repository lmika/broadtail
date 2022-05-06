package stormstore

import (
	"context"

	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
)

type RulesStore struct {
	db *storm.DB
}

func NewRulesStore(dm *DBManager) (*RulesStore, error) {
	db, err := dm.Open(settingsDbFilaname)
	if err != nil {
		return nil, err
	}

	return &RulesStore{db: db}, nil
}

func (rs *RulesStore) List(ctx context.Context) (rules []*models.Rule, err error) {
	err = rs.db.All(&rules)
	return rules, err
}

func (rs *RulesStore) Save(ctx context.Context, rule *models.Rule) error {
	return rs.db.Save(rule)
}

func (rs *RulesStore) Get(ctx context.Context, ruleID uuid.UUID) (*models.Rule, error) {
	var rule models.Rule
	if err := rs.db.One("ID", ruleID, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

func (rs *RulesStore) Delete(ctx context.Context, ruleID uuid.UUID) error {
	return rs.db.DeleteStruct(&models.Rule{ID: ruleID})
}
