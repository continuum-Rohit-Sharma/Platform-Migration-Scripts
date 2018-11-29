package http

import (
	"io"

	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
)

type responseSerializerImpl struct{}

func (ser *responseSerializerImpl) Serialize(res *protocol.Response, dst io.Writer) (err error) {
	responseSerializeRaw(res, dst)
	return
}

func (ser *responseSerializerImpl) Deserialize(src io.Reader) (res *protocol.Response, err error) {
	return responseDeserializeRaw(src)
}
