package procParser

// ParserFactoryImpl is the ParserFactory
type ParserFactoryImpl struct{}

// GetParser returns a Parser implementation
func (ParserFactoryImpl) GetParser() Parser {
	parser := new(simpleParser)
	parser.dependencies.modeHandlerFactory = new(modeHandlerFactoryImpl)
	parser.dependencies.procReaderFactory = new(procReaderFactoryImpl)
	return parser
}

type dependencies interface {
	procReaderFactory
}
