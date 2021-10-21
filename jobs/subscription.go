package jobs

import (
	"container/list"
)

type Subscription struct {
	c       chan Update
	elem    *list.Element
	closeFn func()
}

func (s *Subscription) Chan() chan Update {
	return s.c
}

func (s *Subscription) Close() {
	close(s.c)
	s.closeFn()
}

type Update struct {
	Status string
}
