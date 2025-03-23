package collections

import "sync"

// Thread safe generic implementation of heap.
// Insertion order not maintained on equal values.
type Heap[T any] struct {
	mu         sync.RWMutex
	comparator Comparator[T]
	items      []T
}

// Create an empty generic heap using comparator function to determine order.
func NewHeap[T any](comparator Comparator[T]) *Heap[T] {
	return &Heap[T]{
		comparator: comparator,
		items:      make([]T, 0),
	}
}

// Add item to heap.
func (h *Heap[T]) Push(item T) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.items = append(h.items, item)
	h.siftUp(len(h.items) - 1)
}

// Return and remove highest priority item in heap.
func (h *Heap[T]) Pop() (item T, ok bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.items) == 0 {
		return item, false
	}

	item = h.items[0]

	h.swap(0, len(h.items)-1)
	h.items = h.items[:len(h.items)-1]
	h.siftDown(0)

	return item, true
}

// Return highest priority item in heap without removing.
func (h *Heap[T]) Peek() (item T, ok bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.items) == 0 {
		return item, false
	}

	item = h.items[0]

	return item, true
}

// Return true if heap has no items, false otherwise.
func (h *Heap[T]) Empty() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.items) == 0
}

// Swap two items in heap at the given indices.
func (h *Heap[T]) swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

// Move item at index i down the heap to its correct location.
func (h *Heap[T]) siftDown(i int) {
	for {
		leftChild := i*2 + 1
		if leftChild >= len(h.items) {
			return
		}

		priorityChild := leftChild
		rightChild := leftChild + 1
		if rightChild < len(h.items) && h.comparator(h.items[rightChild], h.items[leftChild]) < 0 {
			priorityChild = rightChild
		}

		if h.comparator(h.items[i], h.items[priorityChild]) < 0 {
			return
		}

		h.swap(i, priorityChild)
		i = priorityChild
	}
}

// Move item at index i up the heap to its correct location
func (h *Heap[T]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if h.comparator(h.items[parent], h.items[i]) < 0 {
			return
		}

		h.swap(i, parent)
		i = parent
	}
}
