package leaderElection

type Interface interface {
	BecomeALeader() (peerID int, isLeader bool, err error)
}
