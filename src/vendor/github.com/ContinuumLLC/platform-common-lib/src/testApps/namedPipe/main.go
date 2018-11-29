package main

import (
	"fmt"
	"os"

	"github.com/ContinuumLLC/platform-common-lib/src/namedpipes"
)

func main() {
	server := namedpipes.GetPipeServer()
	l, err := server.CreatePipe(os.Args[1], &namedpipes.PipeConfig{
		InputBufferSize:  2000,
		OutputBufferSize: 2000,
	})
	if err != nil {
		fmt.Printf("Error in Pipe creation %+v", err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("LRP: %s, Pipe accept failed: %+v", os.Args[1], err)
			continue
		}
		fmt.Printf("Connection %+v", conn)
	}
}
