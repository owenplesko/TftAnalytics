package dedupe

type Task interface {
	Do() error
	ID() string
}
