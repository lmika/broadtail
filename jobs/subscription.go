package jobs

import (
	"container/list"
)

type Subscription struct {
	c       chan SubscriptionEvent
	elem    *list.Element
	closeFn func()
}

func (s *Subscription) Chan() chan SubscriptionEvent {
	return s.c
}

func (s *Subscription) Close() {
	close(s.c)
	s.closeFn()
}

type Update struct {
	Status string
}

type SubscriptionEvent interface{}

type UpdateSubscriptionEvent struct {
	Job    *Job
	Update Update
}

type StateTransitionSubscriptionEvent struct {
	Job       *Job
	FromState JobState
	ToState   JobState
}
