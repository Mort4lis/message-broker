package core

import (
	"errors"
)

type QueueRegistry struct {
	queues map[string]*Queue
}

func NewQueueRegistry() *QueueRegistry {
	return &QueueRegistry{
		queues: make(map[string]*Queue),
	}
}

func (r *QueueRegistry) Register(q *Queue) error {
	if _, ok := r.queues[q.Name()]; ok {
		return errors.New("queue with such name already exists")
	}

	r.queues[q.Name()] = q
	return nil
}

func (r *QueueRegistry) GetByName(name string) (*Queue, bool) {
	q, ok := r.queues[name]
	return q, ok
}

func (r *QueueRegistry) ForEach(fn func(q *Queue) bool) {
	for _, q := range r.queues {
		if !fn(q) {
			break
		}
	}
}
