package namedpipes

import (
	"net"
	"time"
)

// PipeConfig contain configuration for the pipe listener. It directly maps to winio.PipeConfig
type PipeConfig struct {
	SecurityDescriptor string
	MessageMode        bool
	InputBufferSize    int32
	OutputBufferSize   int32
}

// ServerPipe is an interface for server named pipe
type ServerPipe interface {
	CreatePipe(pipeName string, config *PipeConfig) (net.Listener, error)
}

// ClientPipe is an interface for client named pipe
type ClientPipe interface {
	DialPipe(path string, timeout *time.Duration) (net.Conn, error)
}
