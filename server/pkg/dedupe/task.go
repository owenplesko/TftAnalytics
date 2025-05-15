package dedupe

type identifyable interface {
	ID() string
}

type Task[T any] interface {
	identifyable
	Run() T
}

type DelayedCompletionTask[T any] interface {
	identifyable
	Run(complete func()) T
}
