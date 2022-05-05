package models

import "github.com/google/uuid"

type RuleCondition struct {
	FeedID uuid.UUID `req:"feedId"`
	Title  string    `req:"title"`
}

type RuleAction struct {
	Download      bool `req:"download,zero"`
	MarkFavourite bool `req:"markFavourite,zero"`
}

type Rule struct {
	ID        uuid.UUID     `storm:"unique"`
	Name      string        `req:"name"`
	Active    bool          `req:"active,zero"`
	Condition RuleCondition `req:"condition"`
	Action    RuleAction    `req:"action"`
}
