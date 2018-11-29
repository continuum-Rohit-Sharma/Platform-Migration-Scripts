package zookeeper

import (
	"strings"

	"github.com/samuel/go-zookeeper/zk"

	"github.com/ContinuumLLC/platform-common-lib/src/distributed/lock"
)

const zkSeparator = "/"

// ZKClient describe interface for zookeeper client
type ZKClient interface {
	Exists(path string) (bool, *zk.Stat, error)
	Get(path string) ([]byte, *zk.Stat, error)
	Children(path string) ([]string, *zk.Stat, error)
	Set(path string, data []byte, version int32) (*zk.Stat, error)
	Delete(path string, version int32) error
	NewLock(path string, acl []zk.ACL) lock.Locker
	CreateRecursive(childPath string, data []byte, flag int32, acl []zk.ACL) (string, error)
}

type zkClient struct {
	conn *zk.Conn
}

func (client *zkClient) Exists(path string) (bool, *zk.Stat, error) {
	return client.conn.Exists(path)
}

func (client *zkClient) Get(path string) ([]byte, *zk.Stat, error) {
	return client.conn.Get(path)
}

func (client *zkClient) Children(path string) ([]string, *zk.Stat, error) {
	return client.conn.Children(path)
}

func (client *zkClient) Set(path string, data []byte, version int32) (*zk.Stat, error) {
	return client.conn.Set(path, data, version)
}

func (client *zkClient) Delete(path string, version int32) error {
	return client.conn.Delete(path, version)
}

func (client *zkClient) NewLock(path string, acl []zk.ACL) lock.Locker {
	return zk.NewLock(client.conn, path, acl)
}

func (client *zkClient) CreateRecursive(childPath string, data []byte, flag int32, acl []zk.ACL) (path string, err error) {
	path, err = client.conn.Create(childPath, data, flag, acl)
	if err != zk.ErrNoNode {
		return path, err
	}

	// Create parent node.
	parts := strings.Split(childPath, zkSeparator)
	// always skip first argument it should be empty string
	for i := range parts[1:] {
		nPath := strings.Join(parts[:i+2], zkSeparator)

		var exists bool
		exists, _, err = client.conn.Exists(nPath)
		if err != nil {
			return path, err
		}

		if exists {
			continue
		}

		// the last one set real data and flag
		if len(parts)-2 == i {
			return client.conn.Create(nPath, data, flag, acl)
		} else {
			path, err = client.conn.Create(nPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		}

		if err != nil && err != zk.ErrNodeExists {
			return path, err
		}
	}

	return path, err
}
