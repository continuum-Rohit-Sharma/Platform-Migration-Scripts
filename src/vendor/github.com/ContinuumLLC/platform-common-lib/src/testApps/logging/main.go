package main

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ContinuumLLC/platform-common-lib/src/exception"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"

	"time"
)

type attributes struct {
	name  string
	value interface{}
}

func main() {
	c := logging.Config{
		MaxFileSizeInMB: 100,
		OldFileToKeep:   5,
		LogFileName:     `log.log`,
		AllowedLogLevel: logging.INFO,
		ServiceName:     "Plugin1",
	}
	f := logging.GetLoggerFactory()
	logger, err := f.Init(c)
	if err != nil {
		fmt.Println(err)
	}
	w := f.GetWriter().Get()
	var errbuf bytes.Buffer
	f.GetWriter().Set(&errbuf)
	logger.LogWithCorrelationf(logging.INFO, "123", "%s, %s", "Lokesh", "Jain", "Testing", "Logs")
	funcNestError()

	testErr()
	c = logging.Config{
		MaxFileSizeInMB: 100,
		OldFileToKeep:   5,
		LogFileName:     `log.log`,
		AllowedLogLevel: logging.INFO,
		ServiceName:     "Agent-Core",
	}
	logger, _ = logging.GetLoggerFactory().Update(c)
	f.GetWriter().Set(w)
	testLogger()
	logger.LogWithCorrelationf(logging.INFO, "123", "############# Plugin Message ################ \n%s", errbuf.String())
	logger.LogWithCorrelationf(logging.INFO, "123", "############# Plugin Message ################")
	testLogger()
}

func funcReturnError() error {
	return errors.New("ComponentError")
}

func funcNestError() error {
	err := funcReturnError()
	if err != nil {
		err = exception.NewWithMap("AppError", err, map[string]interface{}{
			"one":   1,
			"two":   "2",
			"three": 3.0,
		})
	}
	logger := logging.GetLoggerFactory().Get()
	//logger.LogWithCorrelationf(logging.ERROR, "%%s %s", err)
	logger.LogWithCorrelationf(logging.ERROR, "123", "Map Data : %+v, %v, %v", attributes{
		"aaaa", 22,
	}, attributes{
		"bbb", 23,
	}, err)
	return err
}

func testErr() {
	logger := logging.GetLoggerFactory().Get()
	inner := exception.New("innerCode", nil)
	outer := exception.New("outerCode", inner)
	logger.LogWithCorrelationf(logging.DEBUG, "123", "Inner StackTrace: \n%v", inner)
	logger.LogWithCorrelationf(logging.DEBUG, "123", "Outer StackTrace: \n%+v", outer)
	if logger.IsLogLevel(logging.DEBUG) {
		logger.LogWithCorrelationf(logging.DEBUG, "123", "Inner StackTrace: \n%v", inner)
		logger.LogWithCorrelationf(logging.DEBUG, "123", "Outer StackTrace: \n%v", outer)
	}
}

func testLogger() {
	logger := logging.GetLoggerFactory().Get()
	go logger.Log(logging.DEBUG, "Hello World")
	logger.Log(logging.DEBUG, "Hello World Again")
	time.Sleep(time.Second)
	logger.LogWithCorrelation(logging.INFO, "111", "msg", "a", "b", "c")
	logger.LogWithCorrelation(logging.DEBUG, "111", "msg")
	logger.LogWithCorrelation(logging.ERROR, "131", "msg")
	logger.LogWithCorrelation(logging.FATAL, "121", "msg")
	logger.LogWithCorrelation(logging.WARN, "1231", "msg")
}
