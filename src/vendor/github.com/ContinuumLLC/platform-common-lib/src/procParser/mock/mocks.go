// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/ContinuumLLC/platform-common-lib/src/procParser (interfaces: Parser,ParserFactory)

package mock

import (
	procParser "github.com/ContinuumLLC/platform-common-lib/src/procParser"
	gomock "github.com/golang/mock/gomock"
	io "io"
)

// Mock of Parser interface
type MockParser struct {
	ctrl     *gomock.Controller
	recorder *_MockParserRecorder
}

// Recorder for MockParser (not exported)
type _MockParserRecorder struct {
	mock *MockParser
}

func NewMockParser(ctrl *gomock.Controller) *MockParser {
	mock := &MockParser{ctrl: ctrl}
	mock.recorder = &_MockParserRecorder{mock}
	return mock
}

func (_m *MockParser) EXPECT() *_MockParserRecorder {
	return _m.recorder
}

func (_m *MockParser) Parse(_param0 procParser.Config, _param1 io.ReadCloser) (*procParser.Data, error) {
	ret := _m.ctrl.Call(_m, "Parse", _param0, _param1)
	ret0, _ := ret[0].(*procParser.Data)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockParserRecorder) Parse(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Parse", arg0, arg1)
}

// Mock of ParserFactory interface
type MockParserFactory struct {
	ctrl     *gomock.Controller
	recorder *_MockParserFactoryRecorder
}

// Recorder for MockParserFactory (not exported)
type _MockParserFactoryRecorder struct {
	mock *MockParserFactory
}

func NewMockParserFactory(ctrl *gomock.Controller) *MockParserFactory {
	mock := &MockParserFactory{ctrl: ctrl}
	mock.recorder = &_MockParserFactoryRecorder{mock}
	return mock
}

func (_m *MockParserFactory) EXPECT() *_MockParserFactoryRecorder {
	return _m.recorder
}

func (_m *MockParserFactory) GetParser() procParser.Parser {
	ret := _m.ctrl.Call(_m, "GetParser")
	ret0, _ := ret[0].(procParser.Parser)
	return ret0
}

func (_mr *_MockParserFactoryRecorder) GetParser() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetParser")
}
