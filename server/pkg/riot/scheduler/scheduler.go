package scheduler

import (
	"context"
	"time"

	"github.com/owenplesko/TftAnalytics/pkg/collections"
)

type Scheduler struct {
	pq *collections.PriorityQueue[priorityLock]
}

func New(rate time.Duration) *Scheduler {
	lim := &Scheduler{
		pq: collections.NewPriorityQueue(func(a, b priorityLock) int {
			return b.priority - a.priority
		}),
	}

	go lim.dripper(rate)

	return lim
}

func (lim *Scheduler) Wait(ctx context.Context) chan any {
	lockChan := make(chan any)

	lim.pq.Push(priorityLock{
		lockChan: lockChan,
		priority: GetPriority(ctx),
	})

	return lockChan
}

func (lim *Scheduler) dripper(rate time.Duration) {
	ticker := time.NewTicker(rate)

	for range ticker.C {
		if item, ok := lim.pq.Pop(); ok {
			close(item.lockChan)
		}
	}
}
