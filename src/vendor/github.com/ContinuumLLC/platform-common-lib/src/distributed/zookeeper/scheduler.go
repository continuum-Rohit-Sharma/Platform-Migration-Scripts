package zookeeper

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron"

	"github.com/ContinuumLLC/platform-common-lib/src/distributed/scheduler"
)

var (
	schedulerCron *cron.Cron
	schedulerInit = false
	scheduledJobs []scheduler.ScheduledJob
)

// ScheduledJob is a struct defining the actual scheduled job
// Implementing `Job` interface from cron package and `ScheduledJob` from scheduler package
type Job struct {
	Name     string
	Task     string
	Schedule string
}

func (j Job) GetName() string {
	return j.Name
}

func (j Job) GetTask() string {
	return j.Task
}

func (j Job) GetSchedule() string {
	return j.Schedule
}

//Run initial entry point of a Job
func (j Job) Run() {
	Log.LogInfo("Scheduling job `%s` for execution", j.GetName())
	_, err := Queue.CreateItem(nil, j.Task)
	if err != nil {
		Log.LogError("Scheduler. Couldn't run a distributed job %v, err: %v", j.GetName(), err)
	}
}

//InitDistributedScheduler initializes distributed scheduler
func (SchedulerImpl) DistributedScheduler(ctx context.Context, wg *sync.WaitGroup, jobs []scheduler.ScheduledJob, interval int) error {
	scheduledJobs = jobs
	wg.Add(1)

	go func() {
		defer wg.Done()
	MainLoop:
		for {
			select {
			case <-ctx.Done():
				break MainLoop
			case <-time.After(time.Duration(interval) * time.Second):
				sPeerID, leader, err := LeaderElector.BecomeALeader()
				if err != nil {
					log.Println("become leader got error: ", err)
				}
				if leader && !schedulerInit {
					// I'm a new leader
					startScheduler()
				}
				if sPeerID == undefined && schedulerInit {
					stopScheduler()
				}
			}
		}
	}()
	return nil
}

func startScheduler() {
	Log.LogInfo("I'm a new leader. Initializing scheduler...")
	schedulerInit = true
	schedulerCron = cron.New()

	for _, sj := range scheduledJobs {
		job := Job{
			Name:     sj.GetName(),
			Task:     sj.GetTask(),
			Schedule: sj.GetSchedule(),
		}
		err := schedulerCron.AddJob(job.GetSchedule(), job)
		if err != nil {
			Log.LogError("Couldn't add job %v with schedule %v, err: ", job.GetName(), job.GetSchedule(), err)
			continue
		}
	}
	schedulerCron.Start()
}

func stopScheduler() {
	if schedulerInit {
		Log.LogInfo("Stopping scheduler...")
		schedulerInit = false
		schedulerCron.Stop()
	}
}
