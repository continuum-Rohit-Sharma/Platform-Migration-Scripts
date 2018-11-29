package logging

type Entry struct {
	Impl  loggerImpl
	TxID  string
	Layer int
}

func (entry *Entry) Info(args ...interface{}) {
	entry.Impl.log(INFO, entry.Layer, entry.TxID, args...)
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.Impl.logf(INFO, entry.Layer, entry.TxID, format, args...)
}

func (entry *Entry) Error(args ...interface{}) {
	entry.Impl.log(ERROR, entry.Layer, entry.TxID, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.Impl.logf(ERROR, entry.Layer, entry.TxID, format, args...)
}

func (entry *Entry) Warn(args ...interface{}) {
	entry.Impl.log(WARN, entry.Layer, entry.TxID, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.Impl.logf(WARN, entry.Layer, entry.TxID, format, args...)
}

func (entry *Entry) Debug(args ...interface{}) {
	entry.Impl.log(DEBUG, entry.Layer, entry.TxID, args...)
}

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.Impl.logf(DEBUG, entry.Layer, entry.TxID, format, args...)
}

func (entry *Entry) Fatal(args ...interface{}) {
	entry.Impl.log(FATAL, entry.Layer, entry.TxID, args...)
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.Impl.logf(FATAL, entry.Layer, entry.TxID, format, args...)
}
