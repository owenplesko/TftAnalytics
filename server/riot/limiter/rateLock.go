package limiter

type rateLock struct {
	lockChan chan any
	priority int
}

func (rl rateLock) Compare(other rateLock) bool {
	return rl.priority > other.priority
}
