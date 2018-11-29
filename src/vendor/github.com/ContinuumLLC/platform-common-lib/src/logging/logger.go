package logging

import "io"

const (
	separator        = "-"
	defaultCallDepth = 3
)

var logF = &loggerFactoryImpl{}

//Config is a configuration struct for logging
//A static singleton variable Default log level is INFO
//MaxFileSizeInMB A static singleton variable Max log file size in MB
//OldFileToKeep A static singleton variable max number of old files to be kept before cleaning up
type Config struct {
	AllowedLogLevel LogLevel
	MaxFileSizeInMB int64
	OldFileToKeep   int
	LogFileName     string
	CallDepth       int
	ServiceName     string
}

//SetLogLevel a metod to set Loglevel in configuration for from Give String
func (c *Config) SetLogLevel(level string) {
	switch level {
	case "FATAL":
		c.AllowedLogLevel = FATAL
	case "ERROR":
		c.AllowedLogLevel = ERROR
	case "WARN":
		c.AllowedLogLevel = WARN
	case "INFO":
		c.AllowedLogLevel = INFO
	case "DEBUG":
		c.AllowedLogLevel = DEBUG
	case "TRACE":
		c.AllowedLogLevel = TRACE
	case "OFF":
		c.AllowedLogLevel = OFF
	}
}

// GetLoggerFactory returns Logger Factory
func GetLoggerFactory() LoggerFactory {
	return logF
}

//LoggerFactory is the factory for Logger
type LoggerFactory interface {
	Init(config Config) (Logger, error)
	Update(config Config) (Logger, error)
	Get() Logger
	GetWriter() LogWriter
	GetConfig() Config
	IsInitialized() bool
	//New Depricated method, we should be using Get insead of this
	New(string) Logger
}

//goLoggerFactory is GoLogger Factory
type goLoggerFactory interface {
	New() goLogger
}

// Errorer simple interface for error log
type Errorer interface {
	Error(v ...interface{})
}

//Logger is log operations
type Logger interface {
	Errorer
	LogWithTransaction(level LogLevel, transactionID string, v ...interface{})
	LogWithTransactionf(level LogLevel, transactionID string, format string, v ...interface{})
	WithTransaction(value string) *Entry
	WithTransactionAndLayer(value string, layer int) *Entry
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	IsLogLevel(level LogLevel) bool
	Output(calldepth int, subs, s string) error

	//Depricated Methods
	Log(level LogLevel, v ...interface{})
	Logf(level LogLevel, format string, v ...interface{})
	Write(reader io.Reader) (int64, error)
	LogWithCorrelation(level LogLevel, correlationID string, v ...interface{})
	LogWithCorrelationf(level LogLevel, correlationID string, format string, v ...interface{})

	SetLogLevel(level LogLevel)
}

//goLogger is log operations
type goLogger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Flags() int
	Output(calldepth int, s string) error
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Prefix() string
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	SetFlags(flag int)
	SetOutput(w io.Writer)
	SetPrefix(prefix string)
}

//LogWriter manages Writer to Logger
type LogWriter interface {
	SetLogFile(logFilePath string) error
	Set(out io.Writer)
	SetCloser(out io.Writer, closer io.Closer)
	Reset()
	Get() io.Writer
	TruncateOrRotate() error
	Update(config Config) error
}

//LogLevel is to enable or disable log level
type LogLevel int

const (
	//OFF log level
	OFF LogLevel = 0
	//FATAL log level
	FATAL LogLevel = 10
	//ERROR log level
	ERROR LogLevel = 20
	//WARN log level
	WARN LogLevel = 30
	//INFO log level.
	INFO LogLevel = 40
	//DEBUG log level
	DEBUG LogLevel = 50
	//TRACE log level
	TRACE LogLevel = 60
)

var (
	logLevlMap = map[LogLevel]string{
		FATAL: "FATAL ",
		ERROR: "ERROR ",
		WARN:  "WARN  ",
		INFO:  "INFO  ",
		DEBUG: "DEBUG ",
		TRACE: "TRACE ",
	}
)
