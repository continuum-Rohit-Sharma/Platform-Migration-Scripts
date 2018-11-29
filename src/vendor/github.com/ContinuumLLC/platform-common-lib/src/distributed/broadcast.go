package distributed

import (
	"context"
	"sync"
)

// Event type for event
type Event struct {
	Type    string
	Payload interface{}
}

// BroadcastHandler type for event func
type BroadcastHandler func(e *Event)

// Broadcast define methods for broadcast
type Broadcast interface {
	AddHandler(name string, handler BroadcastHandler)
	Listen(ctx context.Context, wg *sync.WaitGroup)
	CreateEvent(e Event) error
}
