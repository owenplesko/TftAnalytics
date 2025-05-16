package dedupe

import (
	"sync"
)

type Deduplicator struct {
	mu      sync.Mutex
	results map[string]any
}

func New() *Deduplicator {
	return &Deduplicator{
		results: make(map[string]any),
	}
}

func Run[T any](dedupe *Deduplicator, task Task[T]) awaitResult[T] {
	res := getOrInitResult[T](dedupe, task)

	go func() {
		val := task.Run()
		removeResult(dedupe, task)
		res.Announce(val)
	}()

	return res
}

func RunDelayedCompletion[T any](dedupe *Deduplicator, task DelayedCompletionTask[T]) awaitResult[T] {
	res := getOrInitResult[T](dedupe, task)

	go func() {
		complete := func() { removeResult(dedupe, task) }
		val := task.Run(complete)
		res.Announce(val)
	}()

	return res

}

func getOrInitResult[T any](dedupe *Deduplicator, task identifyable) *result[T] {
	dedupe.mu.Lock()
	defer dedupe.mu.Unlock()

	id := task.ID()

	r, ok := dedupe.results[id]
	if ok {
		return r.(*result[T])
	}

	res := newResult[T]()
	dedupe.results[id] = res

	return res
}

func removeResult(dedupe *Deduplicator, task identifyable) {
	dedupe.mu.Lock()
	defer dedupe.mu.Unlock()

	delete(dedupe.results, task.ID())
}
