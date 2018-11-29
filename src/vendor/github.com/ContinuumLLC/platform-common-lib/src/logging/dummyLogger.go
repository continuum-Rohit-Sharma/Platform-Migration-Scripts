package logging

import "io"

//DummyLogger a dummy logger implementation creted for the Unit test cases
type dummyLogger struct {
}

//Log Method to write log statements
func (l dummyLogger) LogWithCorrelation(level LogLevel, correlationID string, v ...interface{}) {
}

//Logf method to write formated log statements
func (l dummyLogger) LogWithCorrelationf(level LogLevel, correlationID string, message string, v ...interface{}) {
}

//Log Method to write log statements
func (l dummyLogger) LogWithTransaction(level LogLevel, correlationID string, v ...interface{}) {
}

//Logf method to write formated log statements
func (l dummyLogger) LogWithTransactionf(level LogLevel, correlationID string, message string, v ...interface{}) {
}

//IsLogLevel Method to find log level
func (l dummyLogger) IsLogLevel(level LogLevel) bool {
	return true
}

//Log Method to write log statements
func (l dummyLogger) Log(level LogLevel, v ...interface{}) {
}

//Logf method to write formated log statements
func (l dummyLogger) Logf(level LogLevel, message string, v ...interface{}) {
}

func (l dummyLogger) Write(reader io.Reader) (int64, error) {
	return 0, nil
}

//SetLogLevel Method to set log level
func (l dummyLogger) SetLogLevel(level LogLevel) {
}

//Output method format and output statement
func (l dummyLogger) Output(calldepth int, subs, message string) error {
	return nil
}

// WithTransaction returns Entry with logger and transaction set
func (l dummyLogger) WithTransaction(value string) *Entry {
	return nil
}

// WithTransactionAndLayer returns Entry with logger, transaction and layer sets
func (l dummyLogger) WithTransactionAndLayer(value string, layer int) *Entry {
	return nil
}

//Info Method to write log statements
func (l dummyLogger) Info(v ...interface{}) {
}

//Infof method to write formated log statements
func (l dummyLogger) Infof(message string, v ...interface{}) {
}

//Error Method to write log statements
func (l dummyLogger) Error(v ...interface{}) {
}

//Errorf method to write formated log statements
func (l dummyLogger) Errorf(message string, v ...interface{}) {
}

//Debug Method to write log statements
func (l dummyLogger) Debug(v ...interface{}) {
}

//Debugf method to write formated log statements
func (l dummyLogger) Debugf(message string, v ...interface{}) {
}

//Fatal Method to write log statements
func (l dummyLogger) Fatal(v ...interface{}) {
}

//Fatalf method to write formated log statements
func (l dummyLogger) Fatalf(message string, v ...interface{}) {
}

//Warn Method to write log statements
func (l dummyLogger) Warn(v ...interface{}) {
}

//Warnf method to write formated log statements
func (l dummyLogger) Warnf(message string, v ...interface{}) {
}
