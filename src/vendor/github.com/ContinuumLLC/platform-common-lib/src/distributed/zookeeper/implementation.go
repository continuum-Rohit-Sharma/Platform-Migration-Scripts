package zookeeper

import (
	"fmt"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"

	"github.com/ContinuumLLC/platform-common-lib/src/distributed/leader-election"
	"github.com/ContinuumLLC/platform-common-lib/src/distributed/lock"
	"github.com/ContinuumLLC/platform-common-lib/src/distributed/queue"
	"github.com/ContinuumLLC/platform-common-lib/src/distributed/scheduler"
)

const (
	nodePrefix         = "node-"
	queuePrefix        = "queue-"
	locksNode          = "locks"
	queueNode          = "queue"
	leaderElectionNode = "leader-election"
)

var (
	zookeeperBasePath string
	LeaderElector     leaderElection.Interface = LeaderElectorImpl{}
	Queue             queue.Interface          = QueueImpl{}
	Scheduler         scheduler.Interface      = SchedulerImpl{}
	Client            ZKClient
	Log               Logger
)

type (
	LeaderElectorImpl struct{}
	QueueImpl         struct{}
	SchedulerImpl     struct{}
)

func Init(zookeeperHosts string, basePath string, logger Logger) error {
	hosts := strings.SplitN(zookeeperHosts, ",", -1)
	conn, _, err := zk.Connect(hosts, 10*time.Second)
	if err != nil {
		return err
	}

	Client = &zkClient{conn: conn}

	if len(basePath) < 1 {
		return fmt.Errorf("incorrect base path: %s", basePath)
	}
	zookeeperBasePath = basePath

	if logger == nil {
		return fmt.Errorf("logger cannot be nil")
	}
	Log = logger

	return nil
}

// NewLock is a wrapper for creating new lock
func NewLock(name string) lock.Locker {
	path := zookeeperBasePath + zkSeparator + locksNode + zkSeparator + name
	acl := zk.WorldACL(zk.PermAll)
	return Client.NewLock(path, acl)
}
