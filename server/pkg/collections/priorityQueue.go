package collections

import "sync"

// Thread safe generic implementation of priority queue.
// Insertion order maintained on equal values.
type PriorityQueue[T any] struct {
	mu         sync.RWMutex
	comparator Comparator[T]
	items      []T
}

// Create an empty generic priority queue using comparator function to determine order.
func NewPriorityQueue[T any](comparator Comparator[T]) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		comparator: comparator,
		items:      make([]T, 0),
	}
}

// Add item to priority queue.
func (pq *PriorityQueue[T]) Push(item T) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	i := pq.insertionIndex(item)
	pq.insert(item, i)
}

// Return and remove item with highest priority.
func (pq *PriorityQueue[T]) Pop() (item T, ok bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if len(pq.items) == 0 {
		return item, false
	}

	item = pq.items[0]
	pq.items = pq.items[1:]

	return item, true
}

// Return item with highest priority without removing.
func (pq *PriorityQueue[T]) Peek() (item T, ok bool) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	if len(pq.items) == 0 {
		return item, false
	}

	item = pq.items[0]

	return item, true
}

// Return true if heap has no items, false otherwise.
func (pq *PriorityQueue[T]) Empty() bool {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return len(pq.items) == 0
}

// Finds the first index i where item > arr[i+1].
// Make sure to lock mutex in calling method.
func (pq *PriorityQueue[T]) insertionIndex(item T) int {
	l, r := 0, len(pq.items)
	for l < r {
		pivot := (l + r) / 2
		if pq.comparator(item, pq.items[pivot]) < 0 {
			r = pivot
		} else {
			l = pivot + 1
		}
	}
	return l
}

// Inserts item at given index.
// Make sure to lock mutex in calling method.
func (pq *PriorityQueue[T]) insert(item T, i int) {
	// add capacity of 1 to slice
	pq.items = append(pq.items, item)
	// shift items to make room for insertion
	copy(pq.items[i+1:], pq.items[i:])
	// insert item into slice
	pq.items[i] = item
}
