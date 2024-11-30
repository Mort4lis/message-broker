package core

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrQueueOverflowed        = errors.New("queue is overflowed")
	ErrReachedSubscriberLimit = errors.New("reached the subscriber limit")
)

type Message []byte

type Queue struct {
	msgCh          chan Message
	evCh           chan event
	doneCh         chan struct{}
	name           string
	maxSubscribers int64
	curSubscribers int64
	once           sync.Once
}

func NewQueue(name string, queueSize, maxSubscribers int64) *Queue {
	return &Queue{
		name:           name,
		maxSubscribers: maxSubscribers,
		evCh:           make(chan event),
		doneCh:         make(chan struct{}, 1),
		msgCh:          make(chan Message, queueSize),
	}
}

func (q *Queue) Name() string {
	return q.name
}

func (q *Queue) Append(msg Message) error {
	select {
	case q.msgCh <- msg:
		q.evCh <- event{Type: sendMessageEventType}
		return nil
	default:
		return ErrQueueOverflowed
	}
}

func (q *Queue) Subscribe() (*Consumer, error) {
	for {
		curSubscribers := atomic.LoadInt64(&q.curSubscribers)
		if curSubscribers+1 >= q.maxSubscribers {
			return nil, ErrReachedSubscriberLimit
		}
		if atomic.CompareAndSwapInt64(&q.curSubscribers, curSubscribers, curSubscribers+1) {
			break
		}
	}

	cons := newConsumer(cap(q.msgCh))
	cons.setCloseCallback(func() {
		q.evCh <- event{Type: unsubscribeEventType, Consumer: cons}
	})
	q.evCh <- event{Type: subscribeEventType, Consumer: cons}
	return cons, nil
}

func (q *Queue) unsubscribe(cons *Consumer) {
	curSubscribers := atomic.LoadInt64(&q.curSubscribers)
	for !atomic.CompareAndSwapInt64(&q.curSubscribers, curSubscribers, curSubscribers-1) {
		curSubscribers = atomic.LoadInt64(&q.curSubscribers)
	}
	close(cons.bufCh)
}

func (q *Queue) StartConsume(ctx context.Context) {
	q.once.Do(func() {
		go q.startConsume(ctx)
	})
}

func (q *Queue) startConsume(ctx context.Context) {
	consumers := make(map[string]*Consumer, q.maxSubscribers)
	defer func() { q.doneCh <- struct{}{} }()

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-q.evCh:
			q.handleEvent(consumers, ev)
		}
	}
}

func (q *Queue) handleEvent(consumers map[string]*Consumer, ev event) {
	switch ev.Type {
	case sendMessageEventType:
	case subscribeEventType:
		consumers[ev.Consumer.ID()] = ev.Consumer
	case unsubscribeEventType:
		q.unsubscribe(ev.Consumer)
		delete(consumers, ev.Consumer.ID())
		return
	}

	if len(consumers) == 0 {
		return
	}

	for {
		select {
		case msg := <-q.msgCh:
			for _, cons := range consumers {
				cons.bufCh <- msg
			}
		default:
			return
		}
	}
}

func (q *Queue) Wait() {
	<-q.doneCh
}
