package core

import (
	"sync/atomic"
)

type Subscriber struct {
	active    atomic.Bool
	outcomeCh chan Message
}

func newSubscriber() *Subscriber {
	return &Subscriber{
		outcomeCh: make(chan Message),
	}
}

func (s *Subscriber) IsActive() bool {
	return s.active.Load()
}

func (s *Subscriber) MessageChannel() <-chan Message {
	return s.outcomeCh
}

func (s *Subscriber) close() {
	close(s.outcomeCh)
	s.active.CompareAndSwap(s.active.Load(), false)
}
