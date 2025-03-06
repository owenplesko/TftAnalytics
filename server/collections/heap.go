package collections

import "sync"

// Thread safe generic implementation of max heap.
// Insertion order not maintained on equal values.
type Heap[T Comparable[T]] struct {
	mu    sync.RWMutex
	items []T
}

func NewHeap[T Comparable[T]]() *Heap[T] {
	return &Heap[T]{
		items: make([]T, 0),
	}
}

func (h *Heap[T]) Push(item T) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.items = append(h.items, item)
	h.siftUp(len(h.items) - 1)
}

func (h *Heap[T]) Pop() T {
	h.mu.Lock()
	defer h.mu.Unlock()

	item := h.items[0]

	h.swap(0, len(h.items)-1)
	h.items = h.items[:len(h.items)-1]
	h.siftDown(0)

	return item
}

func (h *Heap[T]) Peek() T {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.items[0]
}

func (h *Heap[T]) Empty() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.items) == 0
}

func (h *Heap[T]) swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *Heap[T]) siftDown(i int) {
	for {
		leftChild := i*2 + 1
		if leftChild >= len(h.items) {
			return
		}

		biggestChild := leftChild
		rightChild := leftChild + 1
		if rightChild < len(h.items) && h.items[rightChild].Compare(h.items[leftChild]) {
			biggestChild = rightChild
		}

		if h.items[i].Compare(h.items[biggestChild]) {
			return
		}

		h.swap(i, biggestChild)
		i = biggestChild
	}
}

func (h *Heap[T]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if h.items[parent].Compare(h.items[i]) {
			return
		}

		h.swap(i, parent)
		i = parent
	}
}
