package zookeeper

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/samuel/go-zookeeper/zk"

	"github.com/ContinuumLLC/platform-common-lib/src/distributed"
)

const (
	pathForListeners = "broadcast/listeners"
	defaultTimeout   = time.Second
)

var (
	// Broadcast instance
	Broadcast distributed.Broadcast
	// ErrZookeeperNotInit zookeeper is not initialised error
	ErrZookeeperNotInit = errors.New("zookeeper is not initialised use zookeeper.Init")
	mutex               = &sync.Mutex{}
)

type broadcastImpl struct {
	instanceID string
	timeout    time.Duration
	handlers   map[string]distributed.BroadcastHandler
}

// InitBroadcast singleton, thread-safe, returns pointer to *Broadcast
// instanceID should be unique for each instance of microservice
func InitBroadcast(instanceID string, timeout time.Duration) (distributed.Broadcast, error) {
	if Client == nil {
		return nil, ErrZookeeperNotInit
	}

	if Broadcast != nil {
		return Broadcast, nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	if Broadcast != nil {
		return Broadcast, nil
	}

	if timeout <= 0 {
		timeout = defaultTimeout
	}

	Broadcast = &broadcastImpl{
		instanceID: instanceID,
		timeout:    timeout,
		handlers:   make(map[string]distributed.BroadcastHandler),
	}

	return Broadcast, nil
}

// AddHandler adding new handler for listening
func (n *broadcastImpl) AddHandler(name string, handler distributed.BroadcastHandler) {
	n.handlers[name] = handler
}

func (n *broadcastImpl) absolutePath() string {
	return n.listenersPath() + zkSeparator + n.instanceID
}

func (n *broadcastImpl) listenersPath() string {
	return zookeeperBasePath + zkSeparator + pathForListeners
}

func (n *broadcastImpl) subscribe() error {
	exists, _, err := Client.Exists(n.absolutePath())
	if err != nil {
		return err
	}

	if exists {
		return nil
	}
	// we are checking our path, in case if it is absent we are trying to create it
	_, err = Client.CreateRecursive(n.absolutePath(), []byte{}, int32(zk.FlagEphemeral), zk.WorldACL(zk.PermAll))
	return err
}

func logInfo(format string, v ...interface{}) {
	if Log != nil {
		Log.LogInfo(format, v...)
	}
}

// Listen listening input events
func (n *broadcastImpl) Listen(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				logInfo("stopped by context")
				return
			case <-time.After(n.timeout):
				if err := n.subscribe(); err != nil {
					logInfo(err.Error())
				}
				items, err := Queue.GetList(n.instanceID)
				if err != nil && err != zk.ErrNoNode {
					logInfo(err.Error())
					continue
				}
				n.process(items)
			}
		}
	}()
}

func (n *broadcastImpl) process(items []string) {
	for _, item := range items {
		data, err := Queue.GetItemData(n.instanceID, item)
		if err != nil {
			logInfo(err.Error())
			continue
		}

		if err := Queue.RemoveItem(n.instanceID, item); err != nil {
			logInfo(err.Error())
		}

		e := new(distributed.Event)
		if err := json.Unmarshal(data, e); err != nil {
			logInfo(err.Error())
			continue
		}

		if handler, ok := n.handlers[e.Type]; ok {
			handler(e)
		}
	}
}

// CreateEvent creating the new event and send to all subscribers
func (n *broadcastImpl) CreateEvent(e distributed.Event) error {
	// getting all subscribers
	subscribers, _, err := Client.Children(n.listenersPath())
	if err != nil {
		return err
	}

	content, err := json.Marshal(e)
	if err != nil {
		return err
	}

	// sending the message to subscribers
	for _, subscriber := range subscribers {
		_, err := Queue.CreateItem(content, subscriber)
		if err != nil {
			return err
		}
	}

	return nil
}
