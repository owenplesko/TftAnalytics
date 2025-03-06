package collections

type Comparable[T any] interface {
	// returns if item is greater than other
	Compare(other T) bool
}
