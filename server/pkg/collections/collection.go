// Package collections offers generic data structures for types that satisfy the Comparable interface.
//
// Available data structures:
//   - PriorityQueue
//   - Heap
package collections

// Compares generic a and b.
//
// Return values:
//   - Negative: sort a before b
//   - 0: a = b
//   - Positive: sort a after b
type Comparator[T any] func(a, b T) int
