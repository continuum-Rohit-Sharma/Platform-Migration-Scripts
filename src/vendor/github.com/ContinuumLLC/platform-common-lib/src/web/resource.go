package web

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ContinuumLLC/platform-common-lib/src/json"
)

const queryGzipField = "gzip"

// RouteConfig stores configuration of a URL Route
type RouteConfig struct {
	URLPathSuffix string
	URLVars       []string
	Res           Resource
}

// Resource interface is an extension of Route.
// Resource handles requests coming to the Route.
type Resource interface {
	Get(context RequestContext)
	Post(context RequestContext)
	Put(context RequestContext)
	Delete(context RequestContext)
	Others(context RequestContext)
}

// Get405 is a base struct for Get HTTP method not implemented.
type Get405 struct{}

// Get method for Get405 struct to return status 405.
func (Get405) Get(context RequestContext) {
	context.GetResponse().WriteHeader(http.StatusMethodNotAllowed)
}

// Post405 is a base struct for Post HTTP method not implemented.
type Post405 struct{}

// Post method for Post405 struct to return status 405.
func (Post405) Post(context RequestContext) {
	context.GetResponse().WriteHeader(http.StatusMethodNotAllowed)
}

// Put405 is a base struct for Put HTTP method not implemented.
type Put405 struct{}

// Put method for Put405 struct to return status 405.
func (Put405) Put(context RequestContext) {
	context.GetResponse().WriteHeader(http.StatusMethodNotAllowed)
}

// Delete405 is a base struct for Delete HTTP method not implemented.
type Delete405 struct{}

// Delete method for Delete405 struct to return status 405.
func (Delete405) Delete(context RequestContext) {
	context.GetResponse().WriteHeader(http.StatusMethodNotAllowed)
}

// Others405 is a base struct for Others HTTP method not implemented.
type Others405 struct{}

// Others method for Others405 struct to return status 405.
func (Others405) Others(context RequestContext) {
	context.GetResponse().WriteHeader(http.StatusMethodNotAllowed)
}

//CommonResource struct
type CommonResource struct {
	serializer   json.SerializerJSON
	deserializer json.DeserializerJSON
}

//Encode encodes response
func (res CommonResource) Encode(rc RequestContext, w http.ResponseWriter, response interface{}, statusCode int) {
	if statusCode != http.StatusOK {
		http.Error(w, response.(string), statusCode)
		return
	}
	enableGzip := strings.Contains(rc.GetRequest().Header.Get("Accept-Encoding"), queryGzipField)
	var writer io.Writer = w //nolint
	if enableGzip {
		gz, err := res.getGzipWriter(w, gzip.DefaultCompression)
		defer gz.Close() //nolint
		writer = gz
		if err != nil {
			writer = w
		}
	}
	err := res.serializer.Write(writer, response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (res CommonResource) getGzipWriter(w http.ResponseWriter, compressLevel int) (*gzip.Writer, error) {
	if compressLevel != gzip.NoCompression {
		w.Header().Set("Content-Encoding", "gzip")
	}
	return gzip.NewWriterLevel(w, compressLevel)
}

//Decode function reads data and deserializes
func (res CommonResource) Decode(object interface{}, rc RequestContext) error {
	data, err := rc.GetData()
	if err != nil {
		return fmt.Errorf("Unable to read data from request %+v", err)
	}

	err = res.deserializer.ReadString(object, string(data))
	if err != nil {
		return fmt.Errorf("Failed to deserializing data %s :: %+v", string(data), err)
	}
	return nil
}

//GetCommonResource returns CommonResource
func GetCommonResource() CommonResource {
	return CommonResource{
		serializer:   json.FactoryJSONImpl{}.GetSerializerJSON(),
		deserializer: json.FactoryJSONImpl{}.GetDeserializerJSON(),
	}
}
