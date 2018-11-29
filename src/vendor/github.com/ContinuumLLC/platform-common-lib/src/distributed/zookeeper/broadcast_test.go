package zookeeper

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/maraino/go-mock"
	"github.com/samuel/go-zookeeper/zk"

	"github.com/ContinuumLLC/platform-common-lib/src/distributed"
)

var ErrMsg = errors.New("msg")

const (
	handlerName = "event"
	payload     = "Payload"
	subscriber  = "subscriber"
)

func TestInitBroadcast(t *testing.T) {
	const instanceID = "instanceID"
	_, originalClient := InitMock()
	defer Restore(originalClient)

	broadCast, err := InitBroadcast(instanceID, 0)
	if err != nil {
		t.Error(err)
	}
	if broadCast == nil {
		t.Error("broadcast can not be <nil>")
	}

	impl, ok := broadCast.(*broadcastImpl)
	if !ok {
		t.Error("wrong type assertion")
	}

	if impl.timeout != defaultTimeout {
		t.Errorf("expected timeout %d, got %d", defaultTimeout, impl.timeout)
	}
}

func TestBroadcast_AddHandler(t *testing.T) {
	Broadcast.AddHandler(handlerName, func(e *distributed.Event) {})

	impl, ok := Broadcast.(*broadcastImpl)
	if !ok {
		t.Error("wrong type assertion")
	}

	if _, ok := impl.handlers[handlerName]; !ok {
		t.Errorf("broadcast handlers should contains %s", handlerName)
	}
}

func TestBroadcast_CreateEvent(t *testing.T) {
	t.Run("ChildrenError", func(t *testing.T) {
		zkMockObj, originalClient := InitMock()
		defer Restore(originalClient)
		zkMockObj.When("Children", mock.Any).Return([]string{}, &zk.Stat{}, ErrMsg)
		if err := Broadcast.CreateEvent(distributed.Event{Type: handlerName, Payload: payload}); err != ErrMsg {
			t.Errorf("got error %s, expected %s", err, ErrMsg)
		}
	})

	t.Run("MarshalError", func(t *testing.T) {
		zkMockObj, originalClient := InitMock()
		defer Restore(originalClient)
		zkMockObj.When("Children", mock.Any).Return([]string{}, &zk.Stat{}, nil)
		if err := Broadcast.CreateEvent(distributed.Event{Type: handlerName, Payload: make(chan int)}); err == nil {
			t.Error("error can not be <nil>")
		}
	})

	t.Run("Subscribers", func(t *testing.T) {
		zkMockObj, originalClient := InitMock()
		defer Restore(originalClient)
		zkMockObj.When("Children", mock.Any).Return([]string{subscriber}, &zk.Stat{}, nil)
		zkMockObj.When("CreateRecursive", mock.Any, mock.Any, mock.Any, mock.Any).Return("", nil)
		if err := Broadcast.CreateEvent(distributed.Event{Type: handlerName}); err != nil {
			t.Error(err)
		}
	})

	t.Run("SubscribersError", func(t *testing.T) {
		zkMockObj, originalClient := InitMock()
		defer Restore(originalClient)
		zkMockObj.When("Children", mock.Any).Return([]string{subscriber}, &zk.Stat{}, nil)
		zkMockObj.When("CreateRecursive", mock.Any, mock.Any, mock.Any, mock.Any).Return("", ErrMsg)
		if err := Broadcast.CreateEvent(distributed.Event{Type: handlerName}); err != ErrMsg {
			t.Errorf("got error %s, expected %s", err, ErrMsg)
		}
	})
}

func TestBroadcast_subscribe(t *testing.T) {
	impl, ok := Broadcast.(*broadcastImpl)
	if !ok {
		t.Error("wrong type assertion")
	}

	t.Run("Subscribe", func(t *testing.T) {
		zkMockObj, originalClient := InitMock()
		defer Restore(originalClient)
		zkMockObj.When("Exists", mock.Any).Return(false, &zk.Stat{}, nil)
		zkMockObj.When("CreateRecursive", mock.Any, mock.Any, mock.Any, mock.Any).Return("", nil)
		if err := impl.subscribe(); err != nil {
			t.Error(err)
		}
	})

	t.Run("SubscribeError", func(t *testing.T) {
		zkMockObj, originalClient := InitMock()
		defer Restore(originalClient)
		zkMockObj.When("Exists", mock.Any).Return(false, &zk.Stat{}, ErrMsg)
		if err := impl.subscribe(); err != ErrMsg {
			t.Errorf("got error %s, expected %s", err, ErrMsg)
		}
	})

	t.Run("SubscribeExists", func(t *testing.T) {
		zkMockObj, originalClient := InitMock()
		defer Restore(originalClient)
		zkMockObj.When("Exists", mock.Any).Return(true, &zk.Stat{}, nil)
		if err := impl.subscribe(); err != nil {
			t.Error(err)
		}
	})

	t.Run("SubscribeCreateRecursiveError", func(t *testing.T) {
		zkMockObj, originalClient := InitMock()
		defer Restore(originalClient)
		zkMockObj.When("Exists", mock.Any).Return(false, &zk.Stat{}, nil)
		zkMockObj.When("CreateRecursive", mock.Any, mock.Any, mock.Any, mock.Any).Return("", ErrMsg)
		if err := impl.subscribe(); err != ErrMsg {
			t.Errorf("got error %s, expected %s", err, ErrMsg)
		}
	})
}

func TestBroadcast_Listen(t *testing.T) {
	Broadcast.AddHandler(handlerName, func(e *distributed.Event) {
		if e.Payload.(string) != payload {
			t.Errorf("got %s, expected %s", e.Payload.(string), payload)
		}
	})

	impl, ok := Broadcast.(*broadcastImpl)
	if !ok {
		t.Error("wrong type assertion")
	}
	impl.timeout = time.Millisecond

	content, _ := json.Marshal(distributed.Event{Type: handlerName, Payload: payload})
	ctx, cancel := context.WithCancel(context.Background())

	zkMockObj, originalClient := InitMock()
	defer Restore(originalClient)
	zkMockObj.When("Exists", mock.Any).Return(false, &zk.Stat{}, nil).Times(2)
	zkMockObj.When("CreateRecursive", mock.Any, mock.Any, mock.Any, mock.Any).Return("", nil)
	zkMockObj.When("Children", mock.Any).Return([]string{subscriber, subscriber}, &zk.Stat{}, nil).Times(2)
	zkMockObj.When("Get", mock.Any).Return(content, &zk.Stat{}, nil).Times(1)
	zkMockObj.When("Get", mock.Any).Return(content, &zk.Stat{}, nil).Times(1)
	zkMockObj.When("Delete", mock.Any, mock.Any).Return(ErrMsg)

	wg := &sync.WaitGroup{}
	Broadcast.Listen(ctx, wg)
	time.Sleep(time.Millisecond * 10)
	cancel()
	wg.Wait()
}
