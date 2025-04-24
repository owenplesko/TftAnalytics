package dedupe

import "sync"

type Deduplicator struct {
	mu    sync.Mutex
	tasks map[string][]chan error
}

func New() *Deduplicator {
	return &Deduplicator{
		tasks: make(map[string][]chan error),
	}
}

func (dedupe *Deduplicator) Do(task Task) chan error {
	dedupe.mu.Lock()
	defer dedupe.mu.Unlock()

	id := task.ID()
	ch := make(chan error)

	if listeners, exists := dedupe.tasks[id]; exists {
		dedupe.tasks[id] = append(listeners, ch)
		return ch
	}

	listeners := []chan error{ch}
	dedupe.tasks[id] = listeners

	go func() {
		err := task.Do()

		dedupe.mu.Lock()
		defer dedupe.mu.Unlock()
		delete(dedupe.tasks, id)

		for _, listener := range listeners {
			listener <- err
		}
	}()

	return ch
}
