//go:generate goversioninfo

package main

import (
	"io"
	"log"
	"os"
	"strings"
)

var (
	action bool
	logger *log.Logger
)

const glogFlags int = log.Ldate | log.Ltime | log.LUTC | log.Lshortfile

// Usage
// jsonMerge <source.file> <delta.file> [<destination.file>]
// if <destination.file> is not passed, it is set to <source.file>
func main() {
	logfile, err := os.OpenFile("jsonMerge.log", os.O_SYNC|os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	logger = log.New(logfile, "", glogFlags)

	logger.Println("Arguments ", os.Args)
	argsLen := len(os.Args)
	if argsLen < 4 || argsLen > 6 {
		logger.Println(`Usage: 
$ jsonMerge merge/copy <source.file> <delta.file> [<destination.file>]

# if <destination.file> is not passed, it is set to <source.file>`)
	}
	arg := strings.ToLower(os.Args[1])
	switch arg {
	case "merge":
		{
			cfg := &Config{
				Source: os.Args[3],
				Delta:  os.Args[4],
			}
			if argsLen > 5 {
				cfg.Destination = os.Args[5]
			}

			if os.Args[2] == "+" {
				logger.Println("Inside Add")
				cfg.Action = "add"
			} else if os.Args[2] == "-" {
				logger.Println("Inside Remove")
				cfg.Action = "remove"
			} else {
				logger.Println("Inside equal")
				cfg.Action = "equal"
			}

			srv := newMergeService()
			logger.Printf("%+v", srv.Merge(cfg))
		}
	case "move":
		{
			removeFile(os.Args[3])
			moveFile(os.Args[2], os.Args[3])
		}

	case "removefile":
		removeFile(os.Args[4])
	case "rename":
		{
			rename(os.Args[2], os.Args[3])
		}
	}

}

func removeFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		logger.Println(err)
	}
}

func moveFile(source, dest string) {
	_, err := os.Stat(source)

	if !os.IsNotExist(err) {
		from, err := os.Open(source)
		if err != nil {
			logger.Println(err)
			return
		}
		defer from.Close()

		to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			logger.Println(err)
			return
		}

		defer to.Close()
		_, err = io.Copy(to, from)
		if err != nil {
			logger.Println(err)
			return
		}
	} else {
		logger.Println("File ", source, " Does not exist")
	}

}

func rename(source, dest string) {
	err := os.Rename(source, dest)

	if err != nil {
		logger.Println(err)
		return
	}
}
