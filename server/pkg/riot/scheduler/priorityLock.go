package scheduler

type priorityLock struct {
	lockChan chan any
	priority int
}

func (rl priorityLock) Compare(other priorityLock) bool {
	return rl.priority > other.priority
}
