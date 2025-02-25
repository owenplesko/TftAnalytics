package limiter

import (
	"TFTAnalyticsServer/collections"
	"time"
)

type Limiter struct {
	pq *collections.PriorityQueue[requestLock]
}

func New(rate time.Duration) *Limiter {
	lim := &Limiter{
		pq: collections.NewPriorityQueue[requestLock](),
	}

	go lim.dripper(rate)

	return lim
}

func (lim *Limiter) Wait(priority int) chan any {
	lockChan := make(chan any)

	lim.pq.Push(requestLock{
		lockChan: lockChan,
		priority: priority,
	})

	return lockChan
}

func (lim *Limiter) dripper(rate time.Duration) {
	ticker := time.NewTicker(rate)

	for range ticker.C {
		if item, ok := lim.pq.Pop(); ok {
			close(item.lockChan)
		}
	}
}
