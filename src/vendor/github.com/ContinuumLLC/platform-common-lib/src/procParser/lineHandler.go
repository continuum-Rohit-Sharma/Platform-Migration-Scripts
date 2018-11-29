package procParser

import (
	exc "github.com/ContinuumLLC/platform-common-lib/src/exception"

	"strings"
)

const (
	ErrIndexOutOfRange = "Index out of range"
)

//Sepereator characters for different handler
const (
	TablularSeperator string = " "
	KeyValueSeperator string = ":"
)

type modeHandlerFactory interface {
	GetModeHandler(mode Mode) modeHandler
}

type modeHandler interface {
	HandleLine(data string, cfg Config) *Line
}

type modeHandlerFactoryImpl struct {
}

func (modeHandlerFactoryImpl) GetModeHandler(mode Mode) modeHandler {

	switch mode {
	case ModeTabular:
		return new(ModeTabularHandler)
	case ModeKeyValue:
		return new(ModeKeyValueHandler)
	case ModeSeparator:
		return new(ModeSeparatorHandler)
	}
	return nil
}

// ModeSeparatorHandler is to handle the proc data and to split the lines by specified separator
type ModeSeparatorHandler struct {
}

// HandleLine is a handle to split the proc data lines by separator
func (md ModeSeparatorHandler) HandleLine(data string, cfg Config) *Line {
	values := splitLines(data, cfg.Separator)
	return &Line{
		Values: values,
	}
}

//TODO - May be we don't need this handler
type ModeTabularHandler struct {
}

func (md ModeTabularHandler) HandleLine(data string, cfg Config) *Line {
	values := splitLines(data, TablularSeperator)
	line := new(Line)
	line.Values = values
	return line
}

type ModeKeyValueHandler struct {
}

func (md ModeKeyValueHandler) HandleLine(data string, cfg Config) *Line {
	var values []string
	splitValues1 := splitLines(data, KeyValueSeperator)
	for _, v1 := range splitValues1 {
		values1 := splitLines(v1, TablularSeperator)
		for _, v2 := range values1 {
			values = append(values, v2)
		}
	}
	line := new(Line)
	line.Values = values
	return line
}

func splitLines(data string, seperator string) []string {
	var values []string
	data = strings.TrimSpace(data)
	if 0 == len(data) {
		return values
	}
	splitVals := strings.Split(data, seperator)

	for _, v := range splitVals {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		values = append(values, v)
	}
	return values
}

func getKeyValue(values []string, keyIndex int) (string, error) {
	if keyIndex < 0 || keyIndex >= len(values) {
		return "", exc.New(ErrIndexOutOfRange, nil)
	}
	key := values[keyIndex]
	return key, nil
}
