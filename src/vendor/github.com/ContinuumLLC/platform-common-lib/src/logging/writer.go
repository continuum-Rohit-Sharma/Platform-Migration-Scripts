package logging

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"time"

	"sync"

	"github.com/ContinuumLLC/platform-common-lib/src/exception"
)

const oneMB2Byte = 1000000

//LogWriterMap is map to store LogWriter for name
var logWriterMap = make(map[string]LogWriter)

type logWriterImpl struct {
	out     io.Writer
	closer  io.Closer
	current *os.File
	mutex   sync.Mutex
	config  Config
}

//GetLogWriter returns LogWriter for given config and name
func GetLogWriter(config Config, name string) LogWriter {
	if logWriter, ok := logWriterMap[name]; ok {
		return logWriter
	}
	lw := &logWriterImpl{config: config}
	lw.SetLogFile(config.LogFileName)
	logWriterMap[name] = lw
	return lw
}

func (w *logWriterImpl) Set(writer io.Writer) {
	w.SetCloser(writer, nil)
}

func (w *logWriterImpl) SetCloser(writer io.Writer, cl io.Closer) {
	w.out = writer
	w.closer = cl
	w.current = nil
}

func (w *logWriterImpl) Reset() {
	w.Set(getDefaultLogWriter())
}

func (w *logWriterImpl) Get() io.Writer {
	return w.out
}

func (w *logWriterImpl) Update(config Config) error {
	w.config = config
	return nil
}

func (w *logWriterImpl) TruncateOrRotate() error {
	if w.current != nil {
		fileInfo, err := w.current.Stat()
		if err != nil {
			return err
		}
		if fileInfo.Size() > w.config.MaxFileSizeInMB*oneMB2Byte {
			w.mutex.Lock()
			defer w.mutex.Unlock()
			fileInfo, err := w.current.Stat()
			if err != nil {
				return err
			}
			if fileInfo.Size() > w.config.MaxFileSizeInMB*oneMB2Byte {
				err = w.rotate()
			}
			return err
		}
	}
	return nil
}

func (w *logWriterImpl) SetLogFile(logFilePath string) (err error) {
	if logFilePath == "" {
		w.Set(getDefaultLogWriter())
		return
	}
	// check if path exists
	parentDir := filepath.Dir(logFilePath)
	parentDirFileInfo, err := os.Stat(parentDir)
	// path doesn't exist
	if os.IsNotExist(err) {
		os.MkdirAll(parentDir, 0700)
	} else if !parentDirFileInfo.IsDir() {
		return exception.New("LoggingFileParentPathNotDir", err)
	}

	currOut := w.out
	newOut, err := os.OpenFile(logFilePath, os.O_SYNC|os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		w.SetCloser(currOut, w.closer)
		return exception.New("LoggingFileNotOpened", err)
	}

	if w.closer != nil {
		w.closer.Close()
	}

	w.SetCloser(newOut, newOut)
	w.current = newOut
	return
}

func getDefaultLogWriter() io.Writer {
	return os.Stdout
}

func (w *logWriterImpl) rotate() error {
	sFile, err := os.Open(w.current.Name())
	if err != nil {
		return err
	}
	defer sFile.Close()

	dFile, err := os.Create(w.current.Name() + ".old." + time.Now().Format("20060102150405"))
	if err != nil {
		return err
	}
	defer dFile.Close()

	s, err := io.Copy(dFile, sFile) // first var shows number of bytes
	if err != nil {
		return err
	}

	fileInfo, err := w.current.Stat()
	if err != nil {
		return err
	}
	//fmt.Println("Rotating...", w.current.Name(), "Current : ", fileInfo.Size(), "Copied : ", s)
	//w.current.Truncate(fileInfo.Size() - w.config.MaxFileSizeInMB*oneMB2Byte)
	os.Truncate(w.current.Name(), fileInfo.Size()-s)

	err = w.clean()
	if err != nil {
		return err
	}
	return nil
}

func (w *logWriterImpl) clean() error {
	s := filepath.Dir(w.current.Name())
	d, err := os.Open(s)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	base := filepath.Base(w.current.Name())

	var archNames []string
	for _, n := range names {
		if strings.HasPrefix(n, base) && !strings.EqualFold(n, base) {
			f := d.Name() + string(os.PathSeparator) + n
			archNames = append(archNames, f)
		}
	}

	if len(archNames) <= w.config.OldFileToKeep {
		return nil
	}

	sort.Strings(archNames)
	toDel := archNames[0 : len(archNames)-w.config.OldFileToKeep]
	for _, n := range toDel {
		if err := os.Remove(n); err != nil {
			return err
		}
	}
	return nil
}
