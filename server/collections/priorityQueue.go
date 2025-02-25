package collections

import "sync"

// Thread safe generic implementation of priority queue.
// Insertion order maintained on equal values.
type PriorityQueue[T Comparable[T]] struct {
	mu    sync.Mutex
	items []T
}

// Initialize a new priority queue with items of type T.
func NewPriorityQueue[T Comparable[T]]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		items: make([]T, 0),
	}
}

// Add item to correct location in the queue given its priority.
func (pq *PriorityQueue[T]) Push(item T) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	i := pq.insertionIndex(item)
	pq.insert(item, i)
}

// Remove and return the last added item with the highest priority.
func (pq *PriorityQueue[T]) Pop() (T, bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	var item T
	if len(pq.items) == 0 {
		return item, false
	}

	item, pq.items = pq.items[0], pq.items[1:]
	return item, true
}

// Finds the first index i where item > arr[i+1].
// Make sure to lock mutex in calling method.
func (pq *PriorityQueue[T]) insertionIndex(item T) int {
	l, r := 0, len(pq.items)
	for l < r {
		pivot := (l + r) / 2
		if item.Compare(pq.items[pivot]) {
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
