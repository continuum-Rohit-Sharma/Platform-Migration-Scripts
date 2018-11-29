package namedpipes

import (
	"testing"
)

func TestGetPipeServer(t *testing.T) {
	ps := GetPipeServer()
	_, ok := ps.(pipeServer)
	if !ok {
		t.Error("Invalid cast")
	}
	pipe, err := ps.CreatePipe(`\\.\pipe\testPipe1`, nil)
	defer pipe.Close()
	if err != nil {
		t.Error(err)
	}
}

/*
func TestGetPipeClient(t *testing.T) {
	ps := GetPipeServer()
	pc := GetPipeClient()
	сonfig := PipeConfig{
		MessageMode: false,
	}
	_, ok := pc.(pipeClient)
	if !ok {
		t.Error("Invalid cast")
	}

	pipe, err := ps.CreatePipe(`\\.\pipe\testPipe`, &сonfig)
	defer pipe.Close()
	if err != nil {
		t.Error(err)
	}

	// go func() {
	// 	conn, _ := pipe.Accept()
	// 	defer conn.Close()
	// }()

	timeout := time.Second
	conn, err := pc.DialPipe(`\\.\pipe\testPipe`, &timeout)
	if err != nil {
		t.Error(err)
	}
	conn.Close()
}
*/
