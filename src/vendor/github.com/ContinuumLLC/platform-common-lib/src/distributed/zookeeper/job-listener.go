package zookeeper

import (
	"context"
	"sync"
	"time"

	"github.com/ContinuumLLC/platform-common-lib/src/distributed/scheduler"
)

// DistributedJobListener initializes distributed job listeners
func (SchedulerImpl) DistributedJobListener(ctx context.Context, wg *sync.WaitGroup, jobs []scheduler.DistributedJob, interval int) error {
	for _, job := range jobs {
		if err := setDistributedJob(ctx, wg, job, time.Duration(interval)); err != nil {
			Log.LogError("Couldn't init distributed job: %v", err)
			return err
		}
	}
	return nil
}

func setDistributedJob(ctx context.Context, wg *sync.WaitGroup, job scheduler.DistributedJob, interval time.Duration) error {
	wg.Add(1)
	locker := NewLock(job.GetName())

	go func() {
		defer wg.Done()
	MainLoop:
		for {
			select {
			case <-ctx.Done():
				Log.LogError("Warning!!! Distributed Job Listener received ctx.Done(): %v", ctx)
				break MainLoop
			case <-time.After(interval * time.Second):
				if err := locker.Lock(); err != nil {
					Log.LogError("Distributed Job [%s]. Couldn't make a lock, err: %v", job.GetName(), err)
					continue
				}
				items, err := Queue.GetList(job.GetName())
				if err != nil {
					Log.LogError("Couldn't get notification for distributed job: %v, err: %v", job.GetName(), err)
				} else if len(items) > 0 {
					processQueue(ctx, items, job)
				}

				err = locker.Unlock()
				if err != nil {
					Log.LogError("Couldn't unlock lock for distributed job: %v, err: %v", job.GetName(), err)
					locker = NewLock(job.GetName())
					continue
				}
			}
		}
	}()

	return nil
}

func processQueue(ctx context.Context, items []string, job scheduler.DistributedJob) {
	Log.LogInfo("Distributed Job [%s]. Found %d notification(s). Executing callback", job.GetName(), len(items))
	var itemsData [][]byte

	for _, item := range items {
		itemData, err := Queue.GetItemData(job.GetName(), item)
		if err != nil {
			Log.LogError("Distributed Job [%s]. Couldn't get job data, err: %v", job.GetName(), err)
			continue
		}

		if len(itemData) != 0 {
			itemsData = append(itemsData, itemData)
		}

		err = Queue.RemoveItem(job.GetName(), item)
		if err != nil {
			Log.LogError("Distributed Job [%s]. Couldn't remove notification from queue, err: %v", job.GetName(), err)
		}
	}

	job.Callback(ctx, itemsData)
}
