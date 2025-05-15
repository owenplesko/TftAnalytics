package dedupe

type Task[T any] interface {
	ID() string
	Run() T
}
