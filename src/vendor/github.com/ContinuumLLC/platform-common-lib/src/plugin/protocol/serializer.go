package protocol

import (
	"io"
)

// ResponseSerializer exposes methods for operating protocol.Response objects
type ResponseSerializer interface {
	Serialize(res *Response, dst io.Writer) error
	Deserialize(src io.Reader) (res *Response, err error)
}
