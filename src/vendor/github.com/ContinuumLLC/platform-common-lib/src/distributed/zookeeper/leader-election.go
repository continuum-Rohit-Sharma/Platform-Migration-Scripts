package zookeeper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/samuel/go-zookeeper/zk"
)

const undefined = -1

var peerID = undefined

func (LeaderElectorImpl) BecomeALeader() (int, bool, error) {
	var err error

	if peerID == undefined {
		peerID, err = CreatePeerID()
		if err != nil {
			peerID = undefined
			return peerID, false, fmt.Errorf("leader election: couldn't create Peer ID, err: %v", err)
		}
	}

	peers, err := GetActivePeersIDs()
	if err != nil {
		return peerID, false, fmt.Errorf("leader election: couldn't get active Peers IDs, err: %v", err)
	}

	found := false
	minID := undefined
	for _, i := range peers {
		if i < minID || minID == undefined {
			minID = i
		}
		if i == peerID {
			found = true
		}
	}

	if !found {
		peerID = undefined
		return peerID, false, nil
	}
	if minID == peerID {
		return peerID, true, nil
	}
	return peerID, false, nil
}

//CreatePeerID create new peer and return its weight
func CreatePeerID() (int, error) {
	childPath := getLeaderElectionZkPath() + zkSeparator + nodePrefix

	flag := int32(zk.FlagSequence | zk.FlagEphemeral)
	acl := zk.WorldACL(zk.PermAll)

	path, err := Client.CreateRecursive(childPath, nil, flag, acl)
	if err != nil {
		return -1, err
	}
	idx, err := strconv.Atoi(strings.TrimLeft(path[len(childPath):], `0`))
	return idx, err
}

//GetActivePeersIDs returns list of current peer weights
func GetActivePeersIDs() ([]int, error) {
	var res []int
	child, _, err := Client.Children(getLeaderElectionZkPath())
	if err != nil {
		return res, err
	}

	for _, p := range child {
		i, cErr := strconv.Atoi(strings.TrimLeft(p[len(nodePrefix):], `0`))
		if cErr != nil {
			continue
		}
		res = append(res, i)
	}
	return res, nil
}

func getLeaderElectionZkPath() string {
	path := zookeeperBasePath + zkSeparator + leaderElectionNode
	return path
}
