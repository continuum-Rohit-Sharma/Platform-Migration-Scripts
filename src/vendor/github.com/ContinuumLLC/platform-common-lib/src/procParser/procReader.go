package procParser

import (
	"bufio"
	"errors"
	"io"
)

type procReaderFactory interface {
	GetReader(cfg Config, reader io.ReadCloser) procReader
}

type procReader interface {
	ReadLine() (string, error)
	Close() error
}

//ProcReaderFactoryImpl is a factory for geting proc reader
type procReaderFactoryImpl struct{}

//GetReader is a method for reading proc output
func (procReaderFactoryImpl) GetReader(cfg Config, reader io.ReadCloser) procReader {
	readerScanner := procReaderScanner{}
	readerScanner.initReader(cfg, reader)
	return &readerScanner
}

type procReaderScanner struct {
	scanner *bufio.Scanner
	reader  io.ReadCloser
}

func (prs *procReaderScanner) initReader(cfg Config, reader io.ReadCloser) {
	prs.reader = reader
	prs.scanner = bufio.NewScanner(prs.reader)
}

func (prs *procReaderScanner) ReadLine() (string, error) {
	if prs.scanner.Scan() {
		return prs.scanner.Text(), nil
	}
	return "", errors.New(EOFError)
}

func (prs *procReaderScanner) Close() error {
	return prs.reader.Close()
}
