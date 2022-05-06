package models

import (
	"github.com/google/uuid"
	"github.com/lmika/shellwords"
	"strings"
)

type Rule struct {
	ID        uuid.UUID     `storm:"unique"`
	Name      string        `req:"name"`
	Active    bool          `req:"active,zero"`
	Condition RuleCondition `req:"condition"`
	Action    RuleAction    `req:"action"`
}

type RuleWithDescription struct {
	Rule
	Description string
}

type RuleAction struct {
	Download      bool `req:"download,zero"`
	MarkFavourite bool `req:"markFavourite,zero"`
}

func (a RuleAction) Combine(b RuleAction) RuleAction {
	return RuleAction{
		Download:      a.Download || b.Download,
		MarkFavourite: a.MarkFavourite || b.MarkFavourite,
	}
}

type RuleCondition struct {
	FeedID      uuid.UUID `req:"feedId"`
	Title       string    `req:"title"`
	Description string    `req:"description"`
}

func (rs RuleCondition) Matches(target RuleTarget) bool {
	var matches = true

	if rs.FeedID != uuid.Nil {
		matches = matches && rs.FeedID == target.FeedID
	}
	if rs.Title != "" {
		matches = matches && rs.stringMatches(rs.Title, target.Title)
	}
	if rs.Description != "" {
		matches = matches && rs.stringMatches(rs.Description, target.Description)
	}

	return matches
}

func (rs RuleCondition) stringMatches(match, test string) bool {
	if match == "" {
		return true
	}

	test = strings.ToLower(test)
	toks := shellwords.Split(strings.ToLower(match))
	for _, tok := range toks {
		if !strings.Contains(test, tok) {
			return false
		}
	}

	return true
}

type RuleTarget struct {
	FeedID      uuid.UUID
	Title       string
	Description string
}
