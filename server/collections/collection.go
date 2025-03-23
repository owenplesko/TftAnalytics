// Package collections offers generic data structures for types that satisfy the Comparable interface.
//
// Available data structures:
//   - PriorityQueue
//   - Heap
package collections

type Comparator[T any] func(a, b T) int
