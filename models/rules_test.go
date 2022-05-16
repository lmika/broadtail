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

	t.Run("long description", func(t *testing.T) {
		rule := models.RuleCondition{
			Title:       `"The first step"`,
			Description: "JHobz Keizaron",
		}

		assert.True(t, rule.Matches(models.RuleTarget{
			FeedID:      uuid.MustParse("e5d2cd63-c911-457e-9d10-fef9ac2c1a27"),
			Title:       "The First Step - Portal 2 Co-op",
			Description: longDescription,
		}))
	})
}

const longDescription = `
The First Step is a weekly show brought to you by GDQ Hotfix and regular hosts JHobz and Keizaron. Speedrunning can be intimidating to get started with, but it doesn't have to be! Learn just how easy it can be by utilizing semi-blind races -- races of games you've played before, but never tried to speedrun -- and (hopefully) find a few laughs along the way. Tune in live on Thursdays at 7pm Eastern at twitch.tv/gamesdonequick to race alongside the hosts! This edition of The First Step aired on May 12th 2022.
 
Keizaron:
Twitter: https://twitter.com/Keizaron​​​​​​​​​​
Twitch: https://www.twitch.tv/Keizaron

Jhobz:
Twitter: https://twitter.com/J_Hobz​​​​​​​​​​
Twitch: https://www.twitch.tv/JHobz296

Whoishyper:
Twitch: https://www.twitch.tv/whoishyper

ItzBytez:
Twitter: https://twitter.com/ItzBytez
Twitch: https://www.twitch.tv/itzbytez

Interested in hosting a speedrunning event on GamesDoneQuick's Twitch channel? Send us an email to hotfix at gamesdonequick dot com with your event details!
`
