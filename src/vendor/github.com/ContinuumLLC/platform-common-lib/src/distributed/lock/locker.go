package lock

// Locker presents distributed lock/unlock
type Locker interface {
	Lock() error
	Unlock() error
}
