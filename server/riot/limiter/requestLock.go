package limiter

type requestLock struct {
	lockChan chan any
	priority int
}

func (rl requestLock) Compare(other requestLock) bool {
	return rl.priority > other.priority
}
