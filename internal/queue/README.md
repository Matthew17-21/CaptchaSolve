# Queue

A thread-safe FIFO (First-In-First-Out) queue implementation in Go using a slice as the underlying data structure.

## Features

- Thread-safe operations using mutex locks
- Optional maximum capacity constraint
- Efficient slice-based implementation
- Basic queue operations: enqueue, dequeue, peek
- Queue management: length check and clear operations

## Usage

### Creating a Queue

Create an unbounded queue:
```go
queue := NewSliceQueue[int]()
```

Create a bounded queue with maximum capacity:
```go
queue := NewSliceQueue[int](100) // Queue with capacity of 100 elements
```

### Basic Operations

**Enqueue** - Add an element to the end of the queue:
```go
err := queue.Enqueue(42)
if err != nil {
    // Handle queue full error
}
```

**Dequeue** - Remove and return the first element:
```go
val, err := queue.Dequeue()
if err != nil {
    // Handle queue empty error
} else {
    fmt.Printf("Dequeued value: %d\n", val)
}
```

**Peek** - View the first element without removing it:
```go
val, err := queue.Peek()
if err != nil {
    // Handle queue empty error
} else {
    fmt.Printf("First value: %d\n", val)
}
```

### Queue Management

**Check Length** - Get the current number of elements:
```go
length := queue.Len()
```

**Clear Queue** - Remove all elements:
```go
queue.Clear()
```

## Error Handling

The queue operations can return the following errors:

- `ErrQueueFull`: Returned when attempting to enqueue into a full bounded queue
- `ErrQueueEmpty`: Returned when attempting to dequeue or peek from an empty queue

## Thread Safety

All operations on SliceQueue are thread-safe. The implementation uses a `sync.RWMutex` to ensure safe concurrent access:

- Read operations (Peek, Len) use RLock
- Write operations (Enqueue, Dequeue, Clear) use Lock

## Performance Considerations

- The underlying slice grows automatically when needed (for unbounded queues)
- Dequeue operations have O(n) time complexity as they require shifting elements
- All other operations have O(1) time complexity
- Memory usage is proportional to the maximum number of elements that have been in the queue

## Limitations

- Only supports `int` values (modify the implementation if other types are needed)
- For bounded queues, capacity cannot be changed after creation
- No shrinking of underlying slice capacity after dequeue operations

## Example

```go
// Create a bounded queue
queue := NewSliceQueue(3)

// Add elements
queue.Enqueue(1)
queue.Enqueue(2)
queue.Enqueue(3)

// Queue is now full
if err := queue.Enqueue(4); err != nil {
    fmt.Println("Queue is full!")
}

// Remove elements
for queue.Len() > 0 {
    val, _ := queue.Dequeue()
    fmt.Printf("Got value: %d\n", val)
}
```
