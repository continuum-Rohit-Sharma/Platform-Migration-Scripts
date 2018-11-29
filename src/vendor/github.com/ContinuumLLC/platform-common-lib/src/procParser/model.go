package procParser

import "io"

// Data is the parsed output of Proc command
type Data struct {
	Lines []Line
	Map   map[string]Line
}

// Line represents each line in Proc Output
type Line struct {
	Values []string
}

// Mode is the mode of parsing the data
type Mode int

const (
	// ModeKeyValue is for proc data in key:value format and then separated by space
	ModeKeyValue Mode = 10
	// ModeTabular is for proc data in columnar format
	ModeTabular Mode = 20
	// ModeSeparator is for proc data in custom format.
	// If ModeSeparator is set then Separator need to be passed
	ModeSeparator Mode = 30
)

// Config is the configuration parameter to the parser
type Config struct {
	ParserMode    Mode
	IgnoreNewLine bool
	KeyField      int
	Separator     string
}

// Parser is the interface for Proc Parser
type Parser interface {
	Parse(cfg Config, reader io.ReadCloser) (*Data, error)
}

// ParserFactory returns a Parser implementation
type ParserFactory interface {
	GetParser() Parser
}

const (
	//EOFError is an End of File Error
	EOFError = "End of File"
)
