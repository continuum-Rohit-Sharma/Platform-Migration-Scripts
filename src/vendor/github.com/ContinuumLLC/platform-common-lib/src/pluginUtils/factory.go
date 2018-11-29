package pluginUtils

import "io"

//PluginIOReader interface returns a Reader interface
type PluginIOReader interface {
	GetReader() io.Reader
}

//PluginIOWriter interface returns a Writer interface
type PluginIOWriter interface {
	GetWriter() io.Writer
}
