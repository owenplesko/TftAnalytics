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
	dedupe.mu.Lock()
	defer dedupe.mu.Unlock()

	id := task.ID()

	r, ok := dedupe.results[id]
	if ok {
		return r.(*result[T])
	}

	res := newResult[T]()
	dedupe.results[id] = res

	go func() {
		val := task.Run()

		dedupe.mu.Lock()
		defer dedupe.mu.Unlock()

		res.Announce(val)
		delete(dedupe.results, id)
	}()

	return res
}
