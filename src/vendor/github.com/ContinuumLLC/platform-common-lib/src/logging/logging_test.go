package logging

import (
	"bytes"
	"os"
	"testing"
)

func TestLogWriter(t *testing.T) {
	GetLoggerFactory().Init(Config{LogFileName: "log.log"})
	lw := GetLoggerFactory().GetWriter()
	stream := lw.Get()
	if stream == nil {
		t.Error("Stream not initialized")
	}

	newStream := &bytes.Buffer{}
	lw.Set(newStream)
	if lw.Get() != newStream {
		t.Error("Stream not set")
	}

	lw.Reset()
	if lw.Get() != os.Stdout {
		t.Error("Stream not reset")
	}
}

//TODO - Need to fix
// func TestGlog(t *testing.T) {
// 	lgr := GetLoggerFactory().New("test")
// 	if lgr.Prefix() != "test" {
// 		t.Error("Logger Prefix Incorrect")
// 	}
// }
