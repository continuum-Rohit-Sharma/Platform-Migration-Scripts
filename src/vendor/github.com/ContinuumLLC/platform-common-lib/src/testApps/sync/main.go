package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/sync"
	"github.com/ContinuumLLC/platform-common-lib/src/sync/zookeeper"
)

func main() {
	arg := os.Args[1]
	c := logging.Config{
		MaxFileSizeInMB: 100,
		OldFileToKeep:   5,
		LogFileName:     `log.log`,
		AllowedLogLevel: logging.TRACE,
		ServiceName:     "Plugin1",
	}

	f := logging.GetLoggerFactory()
	f.Init(c)

	switch arg {
	case "send":
		send()
		break
	case "listen":
		listen()
		break
	}
	time.Sleep(time.Minute)
}

func send() {
	s := zookeeper.Instance(sync.Config{
		Servers:                []string{"localhost:2181"}, //[]string{"172.28.48.173:2181", "172.28.48.78:2181", "172.28.49.107:2181", "172.28.49.135:2181", "172.28.49.45:2181"},
		SessionTimeoutInSecond: 5,
	})
	go func() {
		for i := 0; i < 500; i++ {
			s.Send("/test", "Test-Data-"+strconv.Itoa(i))
			fmt.Println("Creating Data-", i)
			time.Sleep(1 * time.Second)
		}
	}()
}

func listen() {
	s := zookeeper.Instance(sync.Config{
		Servers:                []string{"localhost:2181"}, //[]string{"172.28.48.173:2181", "172.28.48.78:2181", "172.28.49.107:2181", "172.28.49.135:2181", "172.28.49.45:2181"},
		SessionTimeoutInSecond: 5,
	})
	c := make(chan sync.Response, 1)
	go s.Listen("/test", c)
	for {
		r := <-c
		fmt.Println("Data : ", r.Data, "  Error : ", r.Error)
	}
}
