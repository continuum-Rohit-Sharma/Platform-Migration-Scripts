package zookeeper

import "github.com/samuel/go-zookeeper/zk"

//ListQueue returns list of children in the queue
func (QueueImpl) GetList(queueName string) ([]string, error) {
	children, _, err := Client.Children(getQueueZkPath(queueName))
	return children, err
}

//CreateQueueItem create new sequence child node
func (QueueImpl) CreateItem(data []byte, queueName string) (string, error) {
	childPath := getQueueZkPath(queueName) + zkSeparator + queuePrefix
	flag := int32(zk.FlagSequence)
	acl := zk.WorldACL(zk.PermAll)
	return Client.CreateRecursive(childPath, data, flag, acl)
}

// GetItemData gets node data
func (QueueImpl) GetItemData(queueName, itemName string) ([]byte, error) {
	b, _, err := Client.Get(getQueueZkPath(queueName) + zkSeparator + itemName)
	return b, err
}

// RemoveQueueItem drop node
func (QueueImpl) RemoveItem(queueName string, itemName string) error {
	child := getQueueZkPath(queueName) + zkSeparator + itemName
	return Client.Delete(child, 0)
}

func getQueueZkPath(queueName string) string {
	return zookeeperBasePath + zkSeparator + queueNode + zkSeparator + queueName
}
