package scheduler

import (
	"sync"
	"context"
)

type Interface interface {
	DistributedScheduler(ctx context.Context, wg *sync.WaitGroup, jobs []ScheduledJob, interval int) error
	DistributedJobListener(ctx context.Context, wg *sync.WaitGroup, jobs []DistributedJob, interval int) error
}

type ScheduledJob interface {
	GetName() string
	GetTask() string
	GetSchedule() string
}

type DistributedJob interface {
	GetName() string
	Callback(i ...interface{})
}
