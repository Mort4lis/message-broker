package core

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Consumer struct {
	bufCh   chan Message
	closeCb func()
	id      string
	closed  bool
}

func newConsumer(bufSize int) *Consumer {
	return &Consumer{
		id:    uuid.New().String(),
		bufCh: make(chan Message, bufSize),
	}
}

func (c *Consumer) ID() string {
	return c.id
}

func (c *Consumer) ReadMessage(ctx context.Context) (Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg, ok := <-c.bufCh:
		if !ok {
			return nil, errors.New("closed")
		}
		return msg, nil
	}
}

func (c *Consumer) MessageChan() <-chan Message {
	return c.bufCh
}

func (c *Consumer) setCloseCallback(cb func()) {
	c.closeCb = cb
}

func (c *Consumer) Close() {
	if c.closed {
		return
	}
	if c.closeCb != nil {
		c.closeCb()
	}
	c.closed = true
}
