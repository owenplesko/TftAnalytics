package collections

import (
	"testing"
)

func intMinComparator(a, b int) int {
	return a - b
}

func intMaxComparator(a, b int) int {
	return b - a
}

type PushAndPopTestTable[T any] []struct {
	label      string
	collection interface {
		Push(T)
		Pop() (T, bool)
	}
	input    []T
	expected []T
}

func Test_PushAndPop(t *testing.T) {
	testCases := PushAndPopTestTable[int]{
		{
			label:      "min-heap",
			collection: NewHeap(intMinComparator),
			input:      []int{3, 1, 4, 2},
			expected:   []int{1, 2, 3, 4},
		},
		{
			label:      "max-heap",
			collection: NewHeap(intMaxComparator),
			input:      []int{3, 1, 4, 2},
			expected:   []int{4, 3, 2, 1},
		},
		{
			label:      "min-priorityQueue",
			collection: NewPriorityQueue(intMinComparator),
			input:      []int{3, 1, 4, 2},
			expected:   []int{1, 2, 3, 4},
		},
		{
			label:      "max-priorityQueue",
			collection: NewPriorityQueue(intMaxComparator),
			input:      []int{3, 1, 4, 2},
			expected:   []int{4, 3, 2, 1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			// Push items
			for _, item := range testCase.input {
				testCase.collection.Push(item)
			}

			// Pop items and verify order
			for _, expected := range testCase.expected {
				item, ok := testCase.collection.Pop()
				if !ok {
					t.Fatalf("Expected %d, but collection was empty", expected)
				}
				if item != expected {
					t.Errorf("Expected %d, got %d", expected, item)
				}
			}

			// Pop empty collection
			_, ok := testCase.collection.Pop()
			if ok {
				t.Error("Pop should return false if collection is empty")
			}
		})
	}
}

type PeekTestTable[T any] []struct {
	label      string
	collection interface {
		Push(T)
		Peek() (T, bool)
	}
	input    []T
	expected T
}

func Test_Peek(t *testing.T) {
	testCases := PeekTestTable[int]{
		{
			label:      "min-heap",
			collection: NewHeap(intMinComparator),
			input:      []int{3, 1, 4, 2},
			expected:   1,
		},
		{
			label:      "max-heap",
			collection: NewHeap(intMaxComparator),
			input:      []int{3, 1, 4, 2},
			expected:   4,
		},
		{
			label:      "min-priorityQueue",
			collection: NewPriorityQueue(intMinComparator),
			input:      []int{3, 1, 4, 2},
			expected:   1,
		},
		{
			label:      "max-priorityQueue",
			collection: NewPriorityQueue(intMaxComparator),
			input:      []int{3, 1, 4, 2},
			expected:   4,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			// Peek empty
			_, ok := testCase.collection.Peek()
			if ok {
				t.Errorf("Peek on an empty heap should return false")
			}

			// Push items
			for _, item := range testCase.input {
				testCase.collection.Push(item)
			}

			// Peek top item
			for range 3 {
				val, ok := testCase.collection.Peek()
				if !ok {
					t.Fatalf("Expected Peek to return a value, but collection was empty")
				}
				if val != testCase.expected {
					t.Errorf("Expected Peek to return %d, got %d", testCase.expected, val)
				}
			}
		})
	}
}

type EmptyTestTable[T any] []struct {
	label      string
	collection interface {
		Push(T)
		Pop() (T, bool)
		Empty() bool
	}
}

func Test_Empty(t *testing.T) {
	testCases := EmptyTestTable[int]{
		{
			label:      "heap",
			collection: NewHeap(intMinComparator),
		},
		{
			label:      "priorityQueue",
			collection: NewPriorityQueue(intMinComparator),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			if !testCase.collection.Empty() {
				t.Errorf("New collection should be empty")
			}

			testCase.collection.Push(42)

			if testCase.collection.Empty() {
				t.Errorf("Collection should not be empty after adding an item")
			}

			testCase.collection.Pop()

			if !testCase.collection.Empty() {
				t.Errorf("Collection should be empty after removing the only item")
			}
		})
	}
}

type ItemWithInsertion struct {
	Val       int
	Insertion int
}

func Test_InsertionOrder(t *testing.T) {
	pq := NewPriorityQueue(func(a, b ItemWithInsertion) int {
		return a.Val - b.Val
	})

	// Push items
	for i := range 10 {
		pq.Push(ItemWithInsertion{Val: 42, Insertion: i})
	}

	// Pop items and check insertion
	for i := range 10 {
		item, _ := pq.Pop()

		if item.Insertion != i {
			t.Errorf("Expected item with insertion %d, got %d", i, item.Insertion)
		}
	}
}
