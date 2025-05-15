package dedupe

type result[T any] struct {
	ch  chan any
	val T
}

type awaitResult[T any] interface {
	Await() T
}

func newResult[T any]() *result[T] {
	return &result[T]{
		ch: make(chan any),
	}
}

func (res *result[T]) Await() T {
	<-res.ch
	return res.val
}

func (res *result[T]) Announce(val T) {
	res.val = val
	close(res.ch)
}
