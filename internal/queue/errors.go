package queue

import "errors"

var (
	ErrQueueEmpty = errors.New("queue is empty")
	ErrQueueFull  = errors.New("queue is full")
)
