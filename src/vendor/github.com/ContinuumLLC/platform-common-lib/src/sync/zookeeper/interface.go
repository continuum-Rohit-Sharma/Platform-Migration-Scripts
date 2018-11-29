package zookeeper

import "github.com/samuel/go-zookeeper/zk"

//Connection is a wrapper interface on top of zk.Conn struct to have unit test cases in place
type Connection interface {
	AddAuth(scheme string, auth []byte) error
	Children(path string) ([]string, *zk.Stat, error)
	ChildrenW(path string) ([]string, *zk.Stat, <-chan zk.Event, error)
	Close()
	Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error)
	CreateProtectedEphemeralSequential(path string, data []byte, acl []zk.ACL) (string, error)
	Delete(path string, version int32) error
	Exists(path string) (bool, *zk.Stat, error)
	ExistsW(path string) (bool, *zk.Stat, <-chan zk.Event, error)
	Get(path string) ([]byte, *zk.Stat, error)
	GetACL(path string) ([]zk.ACL, *zk.Stat, error)
	GetW(path string) ([]byte, *zk.Stat, <-chan zk.Event, error)
	Multi(ops ...interface{}) ([]zk.MultiResponse, error)
	Server() string
	SessionID() int64
	Set(path string, data []byte, version int32) (*zk.Stat, error)
	SetACL(path string, acl []zk.ACL, version int32) (*zk.Stat, error)
	SetLogger(l zk.Logger)
	State() zk.State
	Sync(path string) (string, error)
}
