package procParser

import "io"

//ParserDependencies are the dependency structures for the Parser
type ParserDependencies struct {
	modeHandlerFactory
	procReaderFactory
}

type simpleParser struct {
	dependencies ParserDependencies
}

func (parser *simpleParser) Parse(cfg Config, reader io.ReadCloser) (*Data, error) {
	modeHandler := parser.dependencies.modeHandlerFactory.GetModeHandler(cfg.ParserMode)
	procReader := parser.dependencies.procReaderFactory.GetReader(cfg, reader)
	procData := Data{Map: make(map[string]Line)}

	for {
		strLine, err := procReader.ReadLine()
		if err != nil {
			if err.Error() == EOFError {
				break
			} else {
				return nil, err
			}
		}
		if cfg.IgnoreNewLine {
			if strLine == "" {
				continue
			}
		}
		line := modeHandler.HandleLine(strLine, cfg)
		if len(line.Values) > 0 {
			key, err := getKeyValue(line.Values, cfg.KeyField)
			if nil != err {
				return nil, err
			}
			procData.Map[key] = *line
		}
		procData.Lines = append(procData.Lines, *line)
	}
	return &procData, nil
}
