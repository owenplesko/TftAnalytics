package scheduler

import "context"

type contextKey string

const contextPriorityKey contextKey = "priority"

func WithPriority(ctx context.Context, priority int) context.Context {
	return context.WithValue(ctx, contextPriorityKey, priority)
}

func GetPriority(ctx context.Context) int {
	priority, ok := ctx.Value(contextPriorityKey).(int)
	if !ok {
		return 0
	}

	return priority
}
