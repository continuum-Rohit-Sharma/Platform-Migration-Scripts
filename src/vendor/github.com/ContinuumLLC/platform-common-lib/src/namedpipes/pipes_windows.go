// +build windows

package namedpipes

import (
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/ContinuumLLC/platform-common-lib/src/namedpipes/npipe"
	"github.com/ContinuumLLC/platform-common-lib/src/namedpipes/winio"
)

var (
	modkernel32    = syscall.NewLazyDLL("kernel32.dll")
	procCancelIoEx = modkernel32.NewProc("CancelIoEx")

	isWin2003        bool
	checkWin2003Once sync.Once
)

//GetPipeServer creates and returns servePipe implementation
func GetPipeServer() ServerPipe {
	checkWin2003Once.Do(func() {
		errAPI := procCancelIoEx.Find()
		if nil != errAPI {
			isWin2003 = true
		}
	})

	return pipeServer{}
}

type pipeServer struct{}

//CreatePipe takes in a pipename, creates a server named pipe and returns listener for server named pipe
func (pipeServer) CreatePipe(pipeName string, config *PipeConfig) (net.Listener, error) {
	var cfg = &winio.PipeConfig{}
	if config != nil {
		(*cfg).InputBufferSize = (*config).InputBufferSize
		(*cfg).MessageMode = (*config).MessageMode
		(*cfg).OutputBufferSize = (*config).OutputBufferSize
		(*cfg).SecurityDescriptor = (*config).SecurityDescriptor
	}

	if isWin2003 {
		return npipe.Listen(pipeName)
	}
	return winio.ListenPipe(pipeName, cfg)

}

//GetPipeClient creates and returns pipeClient implementation
func GetPipeClient() ClientPipe {
	checkWin2003Once.Do(func() {
		errAPI := procCancelIoEx.Find()
		if nil != errAPI {
			isWin2003 = true
		}
	})

	return pipeClient{}
}

type pipeClient struct{}

//DialPipe takes in a pipename, connects to the pipe and returns connection if successful.
func (pipeClient) DialPipe(path string, timeout *time.Duration) (net.Conn, error) {
	if isWin2003 {
		return npipe.DialTimeout(path, *timeout)
	}
	return winio.DialPipe(path, timeout)
}
