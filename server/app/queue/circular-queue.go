package queue

type CQueueIterator[V any] interface {
	Next() (value *V, ok bool)
}

// CQueue defines a circular queue
type CQueue[V any] struct {
	data     []*V
	capacity int
	head     int
	tail     int
	iterator int
}

// NewCircularQueue creates a queue
func NewCQueue[V any](n int) *CQueue[V] {
	if n == 0 {
		return nil
	}
	return &CQueue[V]{
		data:     make([]*V, n),
		capacity: n,
		head:     0,
		tail:     0,
		iterator: 0,
	}
}

// IsEmpty returns true if queue is empty
func (q *CQueue[V]) IsEmpty() bool {
	return q.head == q.tail
}

// IsFull returns true if queue is full
func (q *CQueue[V]) IsFull() bool {
	return q.head == (q.tail+1)%q.capacity
}

// Enqueue puts a element in the tail of queue
func (q *CQueue[V]) Enqueue(item *V) *CQueue[V] {
	if q.IsFull() {
		panic("tried to enqueue when the queue is full")
	}

	q.data[q.tail] = item
	q.tail = (q.tail + 1) % q.capacity
	return q
}

// Dequeue fetches a element from queue
func (q *CQueue[V]) Dequeue() *V {
	if q.IsEmpty() {
		return nil
	}
	item := q.data[q.head]
	q.head = (q.head + 1) % q.capacity
	q.iterator = q.head // for iteration with Next()
	return item
}
func (q *CQueue[V]) ResetIterator() {
	q.iterator = q.head // reset the iterator
}

// Next increments the iterator
func (q *CQueue[V]) Next() (value *V, ok bool) {
	if q.iterator == q.tail {
		return nil, false
	}
	value = q.data[q.iterator]
	q.iterator = (q.iterator + 1) % q.capacity
	return value, true
}

func (q *CQueue[V]) GetInternalSlice() []*V {
	q.ResetIterator()
	var s []*V = make([]*V, 0, q.capacity)
	for {
		m, ok := q.Next()
		if !ok {
			break
		}
		s = append(s, m)
	}
	return s
}
