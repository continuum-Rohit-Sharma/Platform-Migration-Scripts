package zookeeper

type Logger interface {
	LogInfo(format string, v ...interface{})
	LogError(format string, v ...interface{})
}
