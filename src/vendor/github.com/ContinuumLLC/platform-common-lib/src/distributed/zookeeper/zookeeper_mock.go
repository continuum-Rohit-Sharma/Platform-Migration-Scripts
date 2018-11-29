package zookeeper

import (
	"github.com/maraino/go-mock"
	"github.com/samuel/go-zookeeper/zk"

	"github.com/ContinuumLLC/platform-common-lib/src/distributed/lock"
)

// ClientMock is a mock for ZKClient
type ClientMock struct {
	mock.Mock
}

// InitMock creates new mock for zookeeper and replace original client with it
func InitMock() (mock *ClientMock, original ZKClient) {
	mock = &ClientMock{}
	original = Client
	Client = mock
	return mock, original
}

// Restore set zookeeper client to original one
func Restore(original ZKClient) {
	Client = original
}

// Exists implements ZKClient
func (m *ClientMock) Exists(path string) (bool, *zk.Stat, error) {
	ret := m.Called(path)
	return ret.Bool(0), ret.Get(1).(*zk.Stat), ret.Error(2)
}

// Get implements ZKClient
func (m *ClientMock) Get(path string) ([]byte, *zk.Stat, error) {
	ret := m.Called(path)
	return ret.Get(0).([]byte), ret.Get(1).(*zk.Stat), ret.Error(2)
}

// Children implements ZKClient
func (m *ClientMock) Children(path string) ([]string, *zk.Stat, error) {
	ret := m.Called(path)
	return ret.Get(0).([]string), ret.Get(1).(*zk.Stat), ret.Error(2)
}

// Set implements ZKClient
func (m *ClientMock) Set(path string, data []byte, version int32) (*zk.Stat, error) {
	ret := m.Called(path, data, version)
	return ret.Get(0).(*zk.Stat), ret.Error(1)
}

// Delete implements ZKClient
func (m *ClientMock) Delete(path string, version int32) error {
	return m.Called(path, version).Error(0)
}

// NewLock implements ZKClient
func (m *ClientMock) NewLock(path string, acl []zk.ACL) lock.Locker {
	return m.Called(path, acl).Get(0).(lock.Locker)
}

// CreateRecursive implements ZKClient
func (m *ClientMock) CreateRecursive(childPath string, data []byte, flag int32, acl []zk.ACL) (string, error) {
	ret := m.Called(childPath, data, flag, acl)
	return ret.String(0), ret.Error(1)
}

// LockMock is a mock for ZKLock
type LockMock struct {
	mock.Mock
}

// Lock implements ZKLock
func (m *LockMock) Lock() error {
	return m.Called().Error(0)
}

// Unlock implements ZKLock
func (m *LockMock) Unlock() error {
	return m.Called().Error(0)
}
