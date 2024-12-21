package queue

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSliceQueue(t *testing.T) {
	tests := []struct {
		name        string
		maxCapacity []int
		expectedCap int
	}{
		{
			name:        "unbounded queue",
			expectedCap: 0,
		},
		{
			name:        "bounded queue",
			maxCapacity: []int{5},
			expectedCap: 5,
		},
		{
			name:        "zero capacity becomes unbounded",
			maxCapacity: []int{0},
			expectedCap: 0,
		},
		{
			name:        "negative capacity becomes unbounded",
			maxCapacity: []int{-1},
			expectedCap: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewSliceQueue[int](tt.maxCapacity...)
			require.Equal(t, tt.expectedCap, q.maxCapacity)
		})
	}
}

func TestSliceQueue_Enqueue(t *testing.T) {

	t.Run("unbounded queue", func(t *testing.T) {
		const numElems int = 100
		q := NewSliceQueue[int]()
		for i := 0; i < numElems; i++ {
			if err := q.Enqueue(i); err != nil {
				t.Errorf("Enqueue() error = %v", err)
			}
		}

		require.Equal(t, numElems, q.Len())
	})

	t.Run("bounded queue", func(t *testing.T) {
		q := NewSliceQueue[int](3)
		// Fill queue to capacity
		for i := 0; i < 3; i++ {
			if err := q.Enqueue(i); err != nil {
				t.Errorf("Enqueue() error = %v", err)
			}
		}
		// Attempt to exceed capacity
		if err := q.Enqueue(100); err != ErrQueueFull {
			t.Errorf("Enqueue() error = %v, want %v", err, ErrQueueFull)
		}
	})
}

func TestSliceQueue_Dequeue(t *testing.T) {
	t.Run("empty queue", func(t *testing.T) {
		q := NewSliceQueue[int]()
		_, err := q.Dequeue()
		require.ErrorIs(t, err, ErrQueueEmpty)
	})

	t.Run("FIFO order", func(t *testing.T) {
		q := NewSliceQueue[int]()
		values := []int{1, 2, 3, 4, 5}

		// Enqueue values
		for _, v := range values {
			if err := q.Enqueue(v); err != nil {
				t.Errorf("Enqueue() error = %v", err)
			}
		}

		// Dequeue and verify order
		for _, want := range values {
			got, err := q.Dequeue()
			if err != nil {
				t.Errorf("Dequeue() error = %v", err)
			}

			require.Equal(t, want, got)
		}
	})
}

func TestSliceQueue_Peek(t *testing.T) {
	t.Run("empty queue", func(t *testing.T) {
		q := NewSliceQueue[int]()
		_, err := q.Peek()
		require.ErrorIs(t, err, ErrQueueEmpty)
	})

	t.Run("non-empty queue", func(t *testing.T) {
		q := NewSliceQueue[int]()
		q.Enqueue(42)

		// Peek should return same value multiple times
		for i := 0; i < 3; i++ {
			val, err := q.Peek()
			if err != nil {
				t.Errorf("Peek() error = %v", err)
			}
			if val != 42 {
				t.Errorf("Peek() = %v, want %v", val, 42)
			}
		}

		// Length should remain unchanged
		require.Equal(t, 1, q.Len())
	})
}

func TestSliceQueue_Clear(t *testing.T) {
	q := NewSliceQueue[int]()
	values := []int{1, 2, 3, 4, 5}

	for _, v := range values {
		q.Enqueue(v)
	}

	q.Clear()
	require.Empty(t, q.Len(), "after Clear(), length should be 0")

	_, err := q.Peek()
	require.ErrorIs(t, err, ErrQueueEmpty)
}

func TestSliceQueue_ConcurrentAccess(t *testing.T) {
	q := NewSliceQueue[int]()
	const numGoroutines = 10
	const numOperations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2) // For both enqueuers and dequeuers

	// Launch multiple goroutines that enqueue
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				q.Enqueue(id*numOperations + j)
			}
		}(i)
	}

	// Launch multiple goroutines that dequeue
	successfulDequeues := make(chan int, numGoroutines*numOperations)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				if val, err := q.Dequeue(); err == nil {
					successfulDequeues <- val
				}
			}
		}()
	}

	wg.Wait()
	close(successfulDequeues)

	// Verify that all dequeued values are unique
	seen := make(map[int]bool)
	for val := range successfulDequeues {
		if seen[val] {
			t.Errorf("value %v was dequeued multiple times", val)
		}
		seen[val] = true
	}
}

func TestSliceQueue_Types(t *testing.T) {
	t.Run("string queue", func(t *testing.T) {
		q := NewSliceQueue[string]()
		want := "hello"
		q.Enqueue(want)
		got, err := q.Dequeue()
		if err != nil || got != want {
			t.Errorf("string queue: got %v, want %v, err: %v", got, want, err)
		}
	})

	t.Run("struct queue", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		q := NewSliceQueue[Person]()
		want := Person{Name: "Alice", Age: 30}
		q.Enqueue(want)
		got, err := q.Dequeue()
		if err != nil || got != want {
			t.Errorf("struct queue: got %v, want %v, err: %v", got, want, err)
		}
	})
}
