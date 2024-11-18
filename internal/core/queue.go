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
	maxSubscribers int64
	curSubscribers int64
	msgCh          chan Message
	evCh           chan event
	once           sync.Once
}

func NewQueue(queueSize, maxSubscribers int64) *Queue {
	return &Queue{
		maxSubscribers: maxSubscribers,
		evCh:           make(chan event, 1),
		msgCh:          make(chan Message, queueSize),
	}
}

func (q *Queue) Append(msg Message) error {
	select {
	case q.msgCh <- msg:
		q.evCh <- newSendMessageEvent()
		return nil
	default:
		return ErrQueueOverflowed
	}
}

func (q *Queue) Subscribe() (*Subscriber, error) {
	for {
		curSubscribers := atomic.LoadInt64(&q.curSubscribers)
		if curSubscribers+1 >= q.maxSubscribers {
			return nil, ErrReachedSubscriberLimit
		}
		if atomic.CompareAndSwapInt64(&q.curSubscribers, curSubscribers, curSubscribers+1) {
			break
		}
	}

	sub := newSubscriber()
	q.evCh <- newSubscribeEvent(sub)
	return sub, nil
}

func (q *Queue) Unsubscribe(sub *Subscriber) {
	if !sub.IsActive() {
		return
	}

	curSubscribers := atomic.LoadInt64(&q.curSubscribers)
	for !atomic.CompareAndSwapInt64(&q.curSubscribers, curSubscribers, curSubscribers-1) {
		curSubscribers = atomic.LoadInt64(&q.curSubscribers)
	}
	q.evCh <- newUnsubscribeEvent(sub)
}

func (q *Queue) Run(ctx context.Context) {
	q.once.Do(func() {
		go q.run(ctx)
	})
}

func (q *Queue) run(ctx context.Context) {
	subscribers := make(map[*Subscriber]struct{})

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-q.evCh:
			switch ev.Type {
			case sendMessageEventType:
			case subscribeEventType:
				sub := ev.Meta.(subscriberMeta).Subscriber
				subscribers[sub] = struct{}{}
			case unsubscribeEventType:
				sub := ev.Meta.(subscriberMeta).Subscriber
				if _, ok := subscribers[sub]; ok {
					delete(subscribers, sub)
					sub.close()
				}
				continue
			}

			if len(subscribers) == 0 {
				continue
			}

			msg := <-q.msgCh
			for sub := range subscribers {
				sub.outcomeCh <- msg
			}
		}
	}
}
