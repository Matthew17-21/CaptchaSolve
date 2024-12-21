package tests

import (
	"errors"
	"sync"
	"testing"
)

var (
	ErrQueueEmpty = errors.New("queue is empty")
	ErrQueueFull  = errors.New("queue is full")
)

// SliceQueue implements a generic FIFO queue using a slice.
type SliceQueue[T any] struct {
	data     []T
	capacity int // Optional maximum capacity
	mu       sync.RWMutex
}

// NewSliceQueue creates a new SliceQueue with optional initial capacity.
// If maxCapacity > 0, the queue will be bounded to that size.
func NewSliceQueue[T any](maxCapacity ...int) *SliceQueue[T] {
	var capacity int
	if len(maxCapacity) > 0 && maxCapacity[0] > 0 {
		capacity = maxCapacity[0]
	}
	return &SliceQueue[T]{
		data:     make([]T, 0, capacity),
		capacity: capacity,
	}
}

// Enqueue adds a value to the end of the queue.
// Returns ErrQueueFull if the queue has reached its capacity.
func (q *SliceQueue[T]) Enqueue(val T) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.capacity > 0 && len(q.data) >= q.capacity {
		return ErrQueueFull
	}
	q.data = append(q.data, val)
	return nil
}

// Dequeue removes and returns the first element from the queue.
// Returns ErrQueueEmpty if the queue is empty.
func (q *SliceQueue[T]) Dequeue() (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

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
	q.mu.RLock()
	defer q.mu.RUnlock()

	var zero T
	if len(q.data) == 0 {
		return zero, ErrQueueEmpty
	}
	return q.data[0], nil
}

// Len returns the current number of elements in the queue.
func (q *SliceQueue[T]) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.data)
}

// Clear removes all elements from the queue.
func (q *SliceQueue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data = q.data[:0]
}

// ChannelQueue implements a FIFO queue using channels.
type ChannelQueue struct {
	ch      chan int
	done    chan struct{}
	closed  bool
	mu      sync.RWMutex
	readWg  sync.WaitGroup
	writeWg sync.WaitGroup
}

// NewChannelQueue creates a new channel-based queue with the specified capacity.
func NewChannelQueue(capacity int) *ChannelQueue {
	return &ChannelQueue{
		ch:   make(chan int, capacity),
		done: make(chan struct{}),
	}
}

// Enqueue adds a value to the queue.
// Returns error if the queue is closed.
func (q *ChannelQueue) Enqueue(val int) error {
	q.mu.RLock()
	if q.closed {
		q.mu.RUnlock()
		return errors.New("queue is closed")
	}
	q.writeWg.Add(1)
	q.mu.RUnlock()

	defer q.writeWg.Done()

	select {
	case q.ch <- val:
		return nil
	case <-q.done:
		return errors.New("queue is closed")
	}
}

// Dequeue removes and returns a value from the queue.
// Returns error if the queue is empty or closed.
func (q *ChannelQueue) Dequeue() (int, error) {
	q.mu.RLock()
	if q.closed {
		q.mu.RUnlock()
		return 0, errors.New("queue is closed")
	}
	q.readWg.Add(1)
	q.mu.RUnlock()

	defer q.readWg.Done()

	select {
	case val, ok := <-q.ch:
		if !ok {
			return 0, errors.New("queue is closed")
		}
		return val, nil
	case <-q.done:
		return 0, errors.New("queue is closed")
	}
}

// Close safely closes the queue and waits for all operations to complete.
func (q *ChannelQueue) Close() {
	q.mu.Lock()
	if !q.closed {
		q.closed = true
		close(q.done)
		close(q.ch)
	}
	q.mu.Unlock()

	// Wait for all operations to complete
	q.readWg.Wait()
	q.writeWg.Wait()
}

// Queue represents a thread-safe FIFO queue that can store any type T
type Queue[T any] struct {
	items []T
	mu    sync.RWMutex
}

// New creates a new empty Queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0),
	}
}

// WithCapacity creates a new Queue with initial capacity
func WithCapacity[T any](capacity int) *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0, capacity),
	}
}

// Enqueue adds an item to the end of the queue
func (q *Queue[T]) Enqueue(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

// Dequeue removes and returns the first item in the queue
// Returns error if queue is empty
func (q *Queue[T]) Dequeue() (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var zero T
	if len(q.items) == 0 {
		return zero, errors.New("queue is empty")
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

// Peek returns the first item without removing it
// Returns error if queue is empty
func (q *Queue[T]) Peek() (T, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	var zero T
	if len(q.items) == 0 {
		return zero, errors.New("queue is empty")
	}

	return q.items[0], nil
}

// Length returns the current number of items in the queue
func (q *Queue[T]) Length() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.items)
}

// IsEmpty returns true if the queue has no items
func (q *Queue[T]) IsEmpty() bool {
	return q.Length() == 0
}

// Clear removes all items from the queue
func (q *Queue[T]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = make([]T, 0)
}

// Benchmark functions
func BenchmarkSliceQueue(b *testing.B) {
	const numElements = 1_000_000

	b.Run("Unbounded", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			q := NewSliceQueue[int]()
			b.StartTimer()
			for j := 0; j < numElements; j++ {
				_ = q.Enqueue(j)
			}
			for j := 0; j < numElements; j++ {
				_, _ = q.Dequeue()
			}
			b.StopTimer()
		}
	})

	b.Run("WithCapacity", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			q := NewSliceQueue[int](numElements)
			b.StartTimer()
			for j := 0; j < numElements; j++ {
				_ = q.Enqueue(j)
			}
			for j := 0; j < numElements; j++ {
				_, _ = q.Dequeue()
			}
			b.StopTimer()
		}
	})
}

func BenchmarkChannelQueue(b *testing.B) {
	const numElements = 1_000_000

	b.Run("Buffered", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			q := NewChannelQueue(numElements)
			var wg sync.WaitGroup
			wg.Add(2)

			b.StartTimer()
			// Consumer
			go func() {
				defer wg.Done()
				for j := 0; j < numElements; j++ {
					_, _ = q.Dequeue()
				}
			}()

			// Producer
			go func() {
				defer wg.Done()
				for j := 0; j < numElements; j++ {
					_ = q.Enqueue(j)
				}
				q.Close()
			}()

			wg.Wait()
			b.StopTimer()
		}
	})
}

func BenchmarkFifoQueue(b *testing.B) {
	const numElements = 1_000_000

	b.Run("WithCapacity", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			q := WithCapacity[int](numElements)

			var wg sync.WaitGroup
			wg.Add(2)

			b.StartTimer()
			// Consumer
			go func() {
				defer wg.Done()
				for j := 0; j < numElements; j++ {
					q.Dequeue()
				}
			}()

			// Producer
			go func() {
				defer wg.Done()
				for j := 0; j < numElements; j++ {
					q.Enqueue(j)
				}
				// q.Close()
			}()
			wg.Wait()
			b.StopTimer()
		}
	})
}
