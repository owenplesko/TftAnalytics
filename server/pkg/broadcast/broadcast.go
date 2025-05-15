package broadcast

import "sync"

type Broadcast[T any] struct {
	mu     sync.Mutex
	out    []chan T
	closed bool
}

func New[T any]() *Broadcast[T] {
	return &Broadcast[T]{}
}

func (b *Broadcast[T]) Subscribe(ch chan T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		close(ch)
		return
	}

	b.out = append(b.out, ch)
}

func (b *Broadcast[T]) Unsubscribe(ch chan T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		close(ch)
		return
	}

	// filter ch value out of slice
	n := 0
	for _, x := range b.out {
		if x != ch {
			b.out[n] = x
			n++
		}
	}
	b.out = b.out[:n]
}

func (b *Broadcast[T]) Announce(val T) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, ch := range b.out {
		select {
		case ch <- val:
		default:
		}
	}
}

func (b *Broadcast[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.closed = true

	for _, ch := range b.out {
		close(ch)
	}

	b.out = nil
}
