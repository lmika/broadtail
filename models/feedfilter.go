package models

import "github.com/lmika/shellwords"

type FeedItemFilter struct {
	ContainKeyword []string
	Ordering       string
}

func ParseFeedItemFilter(queryString string) FeedItemFilter {
	tokens := shellwords.Split(queryString)
	return FeedItemFilter{
		ContainKeyword: tokens,
	}
}
