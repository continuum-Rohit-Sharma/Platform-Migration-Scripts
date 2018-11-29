package procParser

import (
	"bytes"
	"io"
	"testing"

	"github.com/ContinuumLLC/platform-common-lib/src/env"
	envMock "github.com/ContinuumLLC/platform-common-lib/src/env/mock"
	"github.com/golang/mock/gomock"
)

type mockReadWriterCloser struct {
	reader io.Reader
}

func (cb *mockReadWriterCloser) Read(p []byte) (int, error) {
	return cb.reader.Read(p)
}

func (cb *mockReadWriterCloser) Close() (err error) {
	return nil
}

func TestReadLine(t *testing.T) {

	ctrl := gomock.NewController(t)
	parser, envMock := SetupMockData(ctrl,
		`MemTotal:        8175496 kB`, false)
	cfg := Config{}
	cfg.IgnoreNewLine = true
	cfg.ParserMode = ModeKeyValue

	reader, err := envMock.GetEnv().GetFileReader("")
	if err != nil {
		t.Errorf("Proc Reader Not available")
		return
	}

	defer reader.Close()
	data, _ := parser.Parse(cfg, reader)

	if len(data.Lines) != 1 {
		t.Errorf("Unexpected number of lines returned")
	}
	if data.Lines[0].Values[0] != "MemTotal" {
		t.Errorf("Unexpected value %s returned, expected value MemTotal", data.Lines[0].Values[0])
	}
	if data.Lines[0].Values[1] != "8175496" {
		t.Errorf("Unexpected value %s returned, expected value MemTotal", data.Lines[0].Values[1])
	}
	if data.Lines[0].Values[2] != "kB" {
		t.Errorf("Unexpected value %s returned, expected value MemTotal", data.Lines[0].Values[2])
	}
}

func TestReadLineEmptyFile(t *testing.T) {

	ctrl := gomock.NewController(t)
	parser, envMock := SetupMockData(ctrl, ``, false)
	cfg := Config{}
	cfg.IgnoreNewLine = true
	cfg.ParserMode = ModeKeyValue

	reader, err := envMock.GetEnv().GetFileReader("")
	if err != nil {
		t.Errorf("Proc Reader Not available")
		return
	}

	defer reader.Close()
	data, _ := parser.Parse(cfg, reader)
	if len(data.Lines) != 0 {
		t.Errorf("Unexpected number of lines returned: %d, expected 0", len(data.Lines))
	}
}

func SetupMockData(ctrl *gomock.Controller, readData string, nilReader bool) (Parser, env.FactoryEnv) {
	data := []byte(readData)
	reader := new(mockReadWriterCloser)
	if !nilReader {
		reader.reader = bytes.NewReader(data)
	}

	mockEnv := envMock.NewMockEnv(ctrl)
	mockEnv.EXPECT().GetFileReader(gomock.Any()).Return(reader, nil)
	mockEnv.EXPECT().GetCommandReader(gomock.Any()).Return(reader, nil)
	mockEnv.EXPECT().GetDirectoryFileCount(gomock.Any()).Return(reader, nil)
	mockFactoryEnv := envMock.NewMockFactoryEnv(ctrl)
	mockFactoryEnv.EXPECT().GetEnv().Return(mockEnv)

	parserFact := new(ParserFactoryImpl)
	parser := parserFact.GetParser()

	// parser := new(simpleParser)
	// parser.dependencies.fileReaderFactory = new(dependenciesImpl)
	// parser.dependencies.modeHandlerFactory = new(modeHandlerFactoryImpl)
	// parser.dependencies.FactoryEnv = mockFactoryEnv

	return parser, mockFactoryEnv
}

func TestParserFactory(t *testing.T) {
	factory := new(ParserFactoryImpl)
	parser := factory.GetParser()
	if parser == nil {
		t.Error("Parser expected, returned nil")
	}
}

// type fileReaderMockImpl struct {
// 	closeError    error
// 	readerError   error
// 	readLineError error
// 	readLineData  string
// 	checkReadLine bool
// 	counter       int
// }

// func (mk *fileReaderMockImpl) GetParser() Parser {
// 	parser := new(simpleParser)
// 	parser.dependencies.fileReaderFactory = mk
// 	parser.dependencies.modeHandlerFactory = new(modeHandlerFactoryImpl)

// 	return parser
// }
// func (mk *fileReaderMockImpl) GetReader(cfg Config) (fileReader, error) {
// 	return mk, mk.readerError
// }
// func (mk *fileReaderMockImpl) ReadLine() (string, error) {
// 	if mk.checkReadLine {
// 		if mk.counter > 1 {
// 			return "", errors.New(EOFError)
// 		}
// 		mk.counter++
// 	}
// 	return mk.readLineData, mk.readLineError
// }

// func (mk *fileReaderMockImpl) close() error {
// 	return mk.closeError
// }

// func TestParseGetReaderError(t *testing.T) {
// 	mock := new(fileReaderMockImpl)
// 	mock.readerError = errors.New("Error getting Reader")
// 	parser := mock.GetParser()
// 	config := new(Config)
// 	_, err := parser.Parse(*config)
// 	if err == nil {
// 		t.Error("Error expected")
// 	}
// }

// func TestParseReadLineError(t *testing.T) {
// 	mock := new(fileReaderMockImpl)
// 	mock.readLineError = errors.New("Non EOF error")
// 	parser := mock.GetParser()
// 	config := new(Config)
// 	config.ParserMode = ModeTabular
// 	_, err := parser.Parse(*config)
// 	if err == nil {
// 		t.Error("Error expected")
// 	}
// }

// func TestParseEOFError(t *testing.T) {
// 	mock := new(fileReaderMockImpl)
// 	mock.readLineError = errors.New(EOFError)
// 	mock.readLineData = "key value"
// 	parser := mock.GetParser()
// 	config := new(Config)
// 	config.ParserMode = ModeTabular
// 	_, err := parser.Parse(*config)
// 	if err != nil {
// 		t.Error("Unexpected Error")
// 	}
// }

// func TestParseIgnoreNewLine(t *testing.T) {
// 	mock := new(fileReaderMockImpl)
// 	mock.readLineData = ""
// 	mock.checkReadLine = true
// 	mock.counter = 1
// 	parser := mock.GetParser()
// 	config := new(Config)
// 	config.IgnoreNewLine = true
// 	config.ParserMode = ModeTabular
// 	_, err := parser.Parse(*config)
// 	if err != nil {
// 		t.Error("Unexpected Error")
// 	}
// }

// func TestParseNoError(t *testing.T) {
// 	mock := new(fileReaderMockImpl)
// 	mock.readLineData = "key value"
// 	mock.checkReadLine = true
// 	mock.counter = 1
// 	parser := mock.GetParser()
// 	config := new(Config)
// 	config.IgnoreNewLine = true
// 	config.ParserMode = ModeTabular
// 	_, err := parser.Parse(*config)
// 	if err != nil {
// 		t.Error("Unexpected Error")
// 	}
// }
