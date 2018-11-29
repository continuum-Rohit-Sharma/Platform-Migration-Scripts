package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ContinuumLLC/platform-common-lib/src/json"
	jMock "github.com/ContinuumLLC/platform-common-lib/src/json/mock"
	"github.com/golang/mock/gomock"
)

func TestCommonResource_Decode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// m := mock.NewMockRequestContext(ctrl)
	// m.EXPECT().GetData().Return([]byte("{}"), errors.New("Error"))
	serializer := jMock.NewMockSerializerJSON(ctrl)
	deserializer := jMock.NewMockDeserializerJSON(ctrl)
	deserializer.EXPECT().ReadString(gomock.Any(), gomock.Any()).Return(errors.New("Error"))
	deserializer.EXPECT().ReadString(gomock.Any(), gomock.Any()).Return(nil)
	type fields struct {
		serializer   json.SerializerJSON
		deserializer json.DeserializerJSON
	}
	type args struct {
		object interface{}
		rc     RequestContext
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{rc: muxRequestContext{request: httptest.NewRequest("POST", "https://abc.com", strings.NewReader("aa"))}},
			fields: fields{
				deserializer: deserializer,
				serializer:   serializer,
			},
			wantErr: true,
		},
		{
			name: "2",
			args: args{rc: muxRequestContext{request: httptest.NewRequest("POST", "https://abc.com", strings.NewReader("aa"))}},
			fields: fields{
				deserializer: deserializer,
				serializer:   serializer,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := CommonResource{
				serializer:   tt.fields.serializer,
				deserializer: tt.fields.deserializer,
			}
			if err := res.Decode(tt.args.object, tt.args.rc); (err != nil) != tt.wantErr {
				t.Errorf("CommonResource.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCommonResource(t *testing.T) {
	_ = GetCommonResource()
}

func getRequest() *http.Request {
	req := httptest.NewRequest("GET", "https://abc.com", strings.NewReader("aa"))
	req.Header.Set("Accept-Encoding", "gzip")
	return req
}

func TestCommonResource_Encode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// m := mock.NewMockRequestContext(ctrl)
	// m.EXPECT().GetData().Return([]byte("{}"), errors.New("Error"))
	serializer := jMock.NewMockSerializerJSON(ctrl)
	deserializer := jMock.NewMockDeserializerJSON(ctrl)
	serializer.EXPECT().Write(gomock.Any(), gomock.Any()).Return(errors.New("Error"))
	serializer.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)
	type fields struct {
		serializer   json.SerializerJSON
		deserializer json.DeserializerJSON
	}
	type args struct {
		rc         RequestContext
		w          http.ResponseWriter
		response   interface{}
		statusCode int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{rc: muxRequestContext{request: httptest.NewRequest("GET", "https://abc.com", strings.NewReader("aa"))},
				w:          httptest.NewRecorder(),
				statusCode: 400,
				response:   "Response",
			},
			fields: fields{
				deserializer: deserializer,
				serializer:   serializer,
			},
		},
		{
			name: "2",
			args: args{
				rc:         muxRequestContext{request: getRequest()},
				w:          httptest.NewRecorder(),
				statusCode: 200,
				response:   "Response",
			},
			fields: fields{
				deserializer: deserializer,
				serializer:   serializer,
			},
		},
		{
			name: "3",
			args: args{rc: muxRequestContext{request: httptest.NewRequest("GET", "https://abc.com", strings.NewReader("aa"))},
				w:          httptest.NewRecorder(),
				statusCode: 200,
				response:   "Response",
			},
			fields: fields{
				deserializer: deserializer,
				serializer:   serializer,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := CommonResource{
				serializer:   tt.fields.serializer,
				deserializer: tt.fields.deserializer,
			}
			res.Encode(tt.args.rc, tt.args.w, tt.args.response, tt.args.statusCode)
		})
	}
}
