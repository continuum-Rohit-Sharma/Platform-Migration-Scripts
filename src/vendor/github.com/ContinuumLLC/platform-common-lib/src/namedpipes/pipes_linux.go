// +build linux darwin

package namedpipes

import (
	"net"
	"os"
	"time"
)

//GetPipeServer creates and returns servePipe implementation
func GetPipeServer() ServerPipe {
	return pipeServer{}
}

type pipeServer struct{}

//CreatePipe takes in a pipename, creates a server named pipe and returns listener for server named pipe
func (pipeServer) CreatePipe(pipeName string, config *PipeConfig) (net.Listener, error) {
	os.Remove(pipeName)
	return net.Listen("unix", pipeName)
}

//GetPipeClient creates and returns pipeClient implementation
func GetPipeClient() ClientPipe {
	return pipeClient{}
}

type pipeClient struct{}

//DialPipe takes in a pipename, connects to the pipe and returns connection if successful.
func (pipeClient) DialPipe(path string, timeout *time.Duration) (net.Conn, error) {
	return net.DialTimeout("unix", path, *timeout)
}
