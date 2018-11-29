package logging

import (
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"
	"time"
)

type loggerFactoryImpl struct {
	w           LogWriter
	config      Config
	initialized bool
}

func (l *loggerFactoryImpl) Init(config Config) (Logger, error) {
	if l.initialized {
		return nil, errors.New("LoggerAlreadyInitialized")
	}
	l.w = &logWriterImpl{config: config}
	err := l.w.SetLogFile(config.LogFileName)
	if err != nil {
		return nil, err
	}
	if config.ServiceName == "" {
		config.ServiceName = "-"
	}
	l.config = config
	l.initialized = true
	return l.Get(), nil
}

func (l *loggerFactoryImpl) Update(config Config) (Logger, error) {
	if !l.initialized {
		return nil, errors.New("LoggerNotInitialized")
	}
	if config.ServiceName == "" {
		config.ServiceName = "-"
	}
	l.w.Update(config)
	l.config = config
	return l.Get(), nil
}

func (l *loggerFactoryImpl) Get() Logger {
	if l.initialized {
		return loggerImpl{
			goLogger: log.New(l.w.Get(), "", glogFlags),
			config:   l.config,
			w:        l.w,
		}
	}
	return dummyLogger{}
}

func (l *loggerFactoryImpl) New(string) Logger {
	return l.Get()
}

func (l *loggerFactoryImpl) GetWriter() LogWriter {
	return l.w
}

func (l *loggerFactoryImpl) GetConfig() Config {
	return l.config
}

func (l *loggerFactoryImpl) IsInitialized() bool {
	return l.initialized
}

type loggerImpl struct {
	goLogger
	config Config
	w      LogWriter

	mu     sync.Mutex // ensures atomic writes; protects the following fields
	prefix string     // prefix to write at beginning of each line
	buf    []byte     // for accumulating text to write
}

func (l loggerImpl) LogWithTransaction(level LogLevel, transactionID string, v ...interface{}) {
	l.log(level, 1, transactionID, v...)
}

func (l loggerImpl) LogWithTransactionf(level LogLevel, transactionID string, message string, v ...interface{}) {
	l.logf(level, 1, transactionID, message, v...)
}

func (l loggerImpl) LogWithCorrelation(level LogLevel, correlationID string, v ...interface{}) {
	l.log(level, 1, correlationID, v...)
}

func (l loggerImpl) LogWithCorrelationf(level LogLevel, correlationID string, message string, v ...interface{}) {
	l.logf(level, 1, correlationID, message, v...)
}

func (l loggerImpl) IsLogLevel(level LogLevel) bool {
	return level <= l.config.AllowedLogLevel
}

func (l loggerImpl) Log(level LogLevel, v ...interface{}) {
	l.log(level, 1, separator, v...)
}

func (l loggerImpl) Logf(level LogLevel, message string, v ...interface{}) {
	l.logf(level, 1, separator, message, v...)
}

func (l loggerImpl) WithTransaction(value string) *Entry {
	return &Entry{
		Impl: l,
		TxID: value,
	}
}

func (l loggerImpl) WithTransactionAndLayer(value string, layer int) *Entry {
	return &Entry{
		Impl:  l,
		TxID:  value,
		Layer: layer,
	}
}

func (l loggerImpl) Info(v ...interface{}) {
	l.log(INFO, 1, separator, v...)
}

func (l loggerImpl) Infof(message string, v ...interface{}) {
	l.logf(INFO, 1, separator, message, v...)
}

func (l loggerImpl) Error(v ...interface{}) {
	l.log(ERROR, 1, separator, v...)
}

func (l loggerImpl) Errorf(message string, v ...interface{}) {
	l.logf(ERROR, 1, separator, message, v...)
}

func (l loggerImpl) Warn(v ...interface{}) {
	l.log(WARN, 1, separator, v...)
}

func (l loggerImpl) Warnf(message string, v ...interface{}) {
	l.logf(WARN, 1, separator, message, v...)
}

func (l loggerImpl) Debug(v ...interface{}) {
	l.log(DEBUG, 1, separator, v...)
}

func (l loggerImpl) Debugf(message string, v ...interface{}) {
	l.logf(DEBUG, 1, separator, message, v...)
}

func (l loggerImpl) Fatal(v ...interface{}) {
	l.log(FATAL, 1, separator, v...)
}

func (l loggerImpl) Fatalf(message string, v ...interface{}) {
	l.logf(FATAL, 1, separator, message, v...)
}

func (l loggerImpl) SetLogLevel(level LogLevel) {
}

func (l loggerImpl) Output(calldepth int, subs, s string) error {
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if glogFlags&(log.Lshortfile|log.Llongfile) != 0 {
		// release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}

	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line, subs)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}

	_, err := l.w.Get().Write(l.buf)
	return err
}

func (l loggerImpl) formatHeader(buf *[]byte, t time.Time, file string, line int, subs string) {
	*buf = append(*buf, l.prefix...)
	if glogFlags&log.LUTC != 0 {
		t = t.UTC()
	}
	if glogFlags&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
		if glogFlags&log.Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if glogFlags&(log.Ltime|log.Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if glogFlags&log.Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}

	*buf = append(*buf, l.config.ServiceName...)
	*buf = append(*buf, ' ')
	*buf = append(*buf, subs...)
	*buf = append(*buf, ' ')

	if glogFlags&(log.Lshortfile|log.Llongfile) != 0 {
		if glogFlags&log.Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

func (l loggerImpl) logf(level LogLevel, layer int, sep string, message string, v ...interface{}) {
	if level <= l.config.AllowedLogLevel {
		if l.config.MaxFileSizeInMB != 0 {
			l.w.TruncateOrRotate()
		}
		val := append([]interface{}{logLevlMap[level]}, v...)
		l.Output(getCallDepth(l.config.CallDepth)+layer, sep, fmt.Sprintf("%s "+message, val...))
	}
}

func (l loggerImpl) log(level LogLevel, layer int, sep string, v ...interface{}) {
	if level <= l.config.AllowedLogLevel {
		if l.config.MaxFileSizeInMB != 0 {
			l.w.TruncateOrRotate()
		}
		val := append([]interface{}{logLevlMap[level]}, v...)
		l.Output(getCallDepth(l.config.CallDepth)+layer, sep, fmt.Sprint(val...))
	}
}

func (l loggerImpl) Write(reader io.Reader) (int64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return io.Copy(l.w.Get(), reader)
}

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func getCallDepth(depth int) int {
	if depth <= 0 {
		return defaultCallDepth
	}
	return depth
}
