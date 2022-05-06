package models_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuleCondition_Matches(t *testing.T) {
	scenarios := []struct {
		condition models.RuleCondition
		expected  bool
	}{
		{condition: models.RuleCondition{FeedID: uuid.MustParse("e5d2cd63-c911-457e-9d10-fef9ac2c1a27")}, expected: true},
		{condition: models.RuleCondition{FeedID: uuid.MustParse("a0f67a1b-ec1f-4590-832e-9c8d8694e40c")}, expected: false},

		{condition: models.RuleCondition{Title: "day triffids"}, expected: true},
		{condition: models.RuleCondition{Title: "triffids \"the day of\""}, expected: true},
		{condition: models.RuleCondition{Title: "none of this"}, expected: false},

		{condition: models.RuleCondition{Description: "film some watching"}, expected: true},
		{condition: models.RuleCondition{Description: "'Some film' watching"}, expected: true},

		{condition: models.RuleCondition{Title: "triffids", Description: "film some"}, expected: true},
		{condition: models.RuleCondition{Title: "triffids", Description: "film fla"}, expected: false},
		{condition: models.RuleCondition{Title: "none of this", Description: "film some"}, expected: false},
	}

	for i, scenario := range scenarios {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.condition.Matches(models.RuleTarget{
				FeedID:      uuid.MustParse("e5d2cd63-c911-457e-9d10-fef9ac2c1a27"),
				Title:       "The day of the triffids",
				Description: "Some film that may be worth watching",
			}))
		})
	}
}
