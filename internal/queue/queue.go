package queue

import "sync"

// SliceQueue implements a generic FIFO queue using a slice.
type SliceQueue[T any] struct {
	mutex       sync.RWMutex
	data        []T
	maxCapacity int // Optional maximum capacity
}

// NewSliceQueue creates a new SliceQueue with optional initial capacity.
// If maxCapacity > 0, the queue will be bounded to that size.
func NewSliceQueue[T any](maxCapacity ...int) *SliceQueue[T] {
	var capacity int
	if len(maxCapacity) > 0 && maxCapacity[0] > 0 {
		capacity = maxCapacity[0]
	}
	return &SliceQueue[T]{
		data:        make([]T, 0, capacity),
		maxCapacity: capacity,
	}
}

// Enqueue adds a value to the end of the queue.
// Returns ErrQueueFull if the queue has reached its capacity.
func (q *SliceQueue[T]) Enqueue(val T) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.maxCapacity > 0 && len(q.data) >= q.maxCapacity {
		return ErrQueueFull
	}
	q.data = append(q.data, val)
	return nil
}

// Dequeue removes and returns the first element from the queue.
// Returns ErrQueueEmpty if the queue is empty.
func (q *SliceQueue[T]) Dequeue() (T, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	var zero T
	if len(q.data) == 0 {
		return zero, ErrQueueEmpty
	}
	val := q.data[0]
	q.data = q.data[1:]
	return val, nil
}

// Peek returns the first element without removing it.
// Returns ErrQueueEmpty if the queue is empty.
func (q *SliceQueue[T]) Peek() (T, error) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	var zero T
	if len(q.data) == 0 {
		return zero, ErrQueueEmpty
	}
	return q.data[0], nil
}

// Len returns the current number of elements in the queue.
func (q *SliceQueue[T]) Len() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return len(q.data)
}

// Clear removes all elements from the queue.
func (q *SliceQueue[T]) Clear() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.data = q.data[:0]
}
