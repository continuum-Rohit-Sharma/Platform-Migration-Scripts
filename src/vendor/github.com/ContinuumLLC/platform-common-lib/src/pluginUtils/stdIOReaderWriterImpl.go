package pluginUtils

import (
	"io"
	"os"
)

//StandardIOReaderImpl is a concreate implementation of IOReaderFactory interface
type StandardIOReaderImpl struct{}

//GetReader is an implementation of interface IOReaderFactory and returns io.Reader
func (StandardIOReaderImpl) GetReader() io.Reader {
	return os.Stdin
}

//StandardIOWriterImpl is a concreate implementation of IOWriterFactory interface
type StandardIOWriterImpl struct{}

//GetWriter is an implementation of interface IOWriterFactory and returns io.Writer
func (StandardIOWriterImpl) GetWriter() io.Writer {
	return os.Stdout
}
