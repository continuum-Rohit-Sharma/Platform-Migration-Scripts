// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/ContinuumLLC/platform-common-lib/src/clar (interfaces: ServiceInit,ServiceInitFactory)

package mock

import (
	clar "github.com/ContinuumLLC/platform-common-lib/src/clar"
	gomock "github.com/golang/mock/gomock"
)

// Mock of ServiceInit interface
type MockServiceInit struct {
	ctrl     *gomock.Controller
	recorder *_MockServiceInitRecorder
}

// Recorder for MockServiceInit (not exported)
type _MockServiceInitRecorder struct {
	mock *MockServiceInit
}

func NewMockServiceInit(ctrl *gomock.Controller) *MockServiceInit {
	mock := &MockServiceInit{ctrl: ctrl}
	mock.recorder = &_MockServiceInitRecorder{mock}
	return mock
}

func (_m *MockServiceInit) EXPECT() *_MockServiceInitRecorder {
	return _m.recorder
}

func (_m *MockServiceInit) GetConfigPath() string {
	ret := _m.ctrl.Call(_m, "GetConfigPath")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockServiceInitRecorder) GetConfigPath() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetConfigPath")
}

func (_m *MockServiceInit) GetLogFilePath() string {
	ret := _m.ctrl.Call(_m, "GetLogFilePath")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockServiceInitRecorder) GetLogFilePath() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetLogFilePath")
}

func (_m *MockServiceInit) SetupOsArgs(_param0 string, _param1 string, _param2 []string, _param3 int, _param4 int) {
	_m.ctrl.Call(_m, "SetupOsArgs", _param0, _param1, _param2, _param3, _param4)
}

func (_mr *_MockServiceInitRecorder) SetupOsArgs(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetupOsArgs", arg0, arg1, arg2, arg3, arg4)
}

// Mock of ServiceInitFactory interface
type MockServiceInitFactory struct {
	ctrl     *gomock.Controller
	recorder *_MockServiceInitFactoryRecorder
}

// Recorder for MockServiceInitFactory (not exported)
type _MockServiceInitFactoryRecorder struct {
	mock *MockServiceInitFactory
}

func NewMockServiceInitFactory(ctrl *gomock.Controller) *MockServiceInitFactory {
	mock := &MockServiceInitFactory{ctrl: ctrl}
	mock.recorder = &_MockServiceInitFactoryRecorder{mock}
	return mock
}

func (_m *MockServiceInitFactory) EXPECT() *_MockServiceInitFactoryRecorder {
	return _m.recorder
}

func (_m *MockServiceInitFactory) GetServiceInit() clar.ServiceInit {
	ret := _m.ctrl.Call(_m, "GetServiceInit")
	ret0, _ := ret[0].(clar.ServiceInit)
	return ret0
}

func (_mr *_MockServiceInitFactoryRecorder) GetServiceInit() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetServiceInit")
}
