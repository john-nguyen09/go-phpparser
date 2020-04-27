/*
Package queue provides a fast, ring-buffer queue based on the version suggested by Dariusz Górecki.
Using this instead of other, simpler, queue implementations (slice+append or linked list) provides
substantial memory and time benefits, and fewer GC pauses.

The queue implemented here is as fast as it is for an additional reason: it is *not* thread-safe.
*/
package queue

// Rune represents a single instance of the queue data structure.
type Rune struct {
	buf               []rune
	head, tail, count int
}

// NewRune constructs and returns a new Queue.
func NewRune() *Rune {
	return &Rune{
		buf: make([]rune, minQueueLen),
	}
}

// Length returns the number of elements currently stored in the queue.
func (q *Rune) Length() int {
	return q.count
}

// resizes the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (q *Rune) resize() {
	newBuf := make([]rune, q.count<<1)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

// Add puts an element on the end of the queue.
func (q *Rune) Add(elem rune) {
	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = elem
	// bitwise modulus
	q.tail = (q.tail + 1) & (len(q.buf) - 1)
	q.count++
}

// Peek returns the element at the head of the queue. This call panics
// if the queue is empty.
func (q *Rune) Peek() rune {
	if q.count <= 0 {
		panic("queue: Peek() called on empty queue")
	}
	return q.buf[q.head]
}

// Get returns the element at index i in the queue. If the index is
// invalid, the call will panic. This method accepts both positive and
// negative index values. Index 0 refers to the first element, and
// index -1 refers to the last.
func (q *Rune) Get(i int) rune {
	// If indexing backwards, convert to positive index.
	if i < 0 {
		i += q.count
	}
	if i < 0 || i >= q.count {
		panic("queue: Get() called with index out of range")
	}
	// bitwise modulus
	return q.buf[(q.head+i)&(len(q.buf)-1)]
}

// Remove removes and returns the element from the front of the queue. If the
// queue is empty, the call will panic.
func (q *Rune) Remove() rune {
	if q.count <= 0 {
		panic("queue: Remove() called on empty queue")
	}
	ret := q.buf[q.head]
	q.buf[q.head] = 0
	// bitwise modulus
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.count--
	// Resize down if buffer 1/4 full.
	if len(q.buf) > minQueueLen && (q.count<<2) == len(q.buf) {
		q.resize()
	}
	return ret
}

// Slice creates a slice from the queue
func (q *Rune) Slice(start int, end int) []rune {
	length := end - start
	start = (q.head + start) & (len(q.buf) - 1)
	end = (q.head + end) & (len(q.buf) - 1)
	newBuf := make([]rune, length)
	if end > start {
		return q.buf[start:end]
	}
	n := copy(newBuf, q.buf[start:])
	copy(newBuf[n:], q.buf[:end])
	return newBuf
}
