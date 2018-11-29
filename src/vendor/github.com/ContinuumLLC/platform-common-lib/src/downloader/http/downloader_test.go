package http

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ContinuumLLC/platform-common-lib/src/downloader"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/utils"
	"github.com/ContinuumLLC/platform-common-lib/src/webClient"
	cMock "github.com/ContinuumLLC/platform-common-lib/src/webClient/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	downloadLocation = "download"
)

func TestGetDownloader(t *testing.T) {
	service := GetDownloader(webClient.BasicClient, webClient.ClientConfig{})
	_, ok := service.(serviceImpl)
	if !ok {
		t.Error("Invalid serviceImpl")
	}
}

func Test_serviceImpl_Download(t *testing.T) {
	ctrl := gomock.NewController(t)

	clientMock := cMock.NewMockClientService(ctrl)
	clientMock.EXPECT().Do(gomock.Any()).Return(nil, errors.New("Error"))

	resp3 := &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader("Read"))}
	clientMock3 := cMock.NewMockClientService(ctrl)
	clientMock3.EXPECT().Do(gomock.Any()).Return(resp3, nil)

	type fields struct {
		client webClient.ClientService
		log    logging.Logger
	}
	type args struct {
		conf downloader.Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "1",
			fields:  fields{client: clientMock, log: logging.GetLoggerFactory().Get()},
			args:    args{conf: downloader.Config{URL: ":::>>>"}}, //Wrong URL
			wantErr: true,
		},
		{
			name:    "2",
			fields:  fields{client: clientMock, log: logging.GetLoggerFactory().Get()},
			args:    args{conf: downloader.Config{URL: "http://test", DownloadLocation: downloadLocation}},
			wantErr: true,
		},
		{
			name:    "3",
			fields:  fields{client: clientMock3, log: logging.GetLoggerFactory().Get()},
			args:    args{conf: downloader.Config{URL: "http://test", DownloadLocation: downloadLocation, FileName: fileName}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serviceImpl{
				client: tt.fields.client,
				log:    tt.fields.log,
			}
			if err := s.Download(&tt.args.conf); (err != nil) != tt.wantErr {
				t.Errorf("serviceImpl.Download() error = %v, wantErr %v", err, tt.wantErr)
			}
			defer func() {
				os.RemoveAll(downloadLocation)
				os.RemoveAll(fileName)
			}()
		})
	}
}

func TestGenerateFileName(t *testing.T) {
	t.Run("EmptyName", func(t *testing.T) {
		conf := downloader.Config{FileName: "test"}
		generateFileName(&conf, &http.Response{})
		assert.NotEqual(t, "", conf.FileName)
	})

	t.Run("GenerateName", func(t *testing.T) {
		conf := downloader.Config{FileName: "", KeepOriginalName: false}
		r := httptest.NewRecorder()
		resp := r.Result()
		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		resp.Request = req
		generateFileName(&conf, resp)
		msg := []byte(resp.Request.URL.String())
		checkSum := utils.GetChecksum(msg)
		assert.Equal(t, checkSum, conf.FileName)
	})

	t.Run("GetExtFromHeader_1", func(t *testing.T) {
		conf := downloader.Config{FileName: "", KeepOriginalName: false}
		r := httptest.NewRecorder()
		resp := r.Result()
		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		resp.Request = req
		resp.Header.Add("content-disposition", "attachment; filename=foo.exe")
		msg := []byte(resp.Request.URL.String())
		checkSum := utils.GetChecksum(msg)
		generateFileName(&conf, resp)
		assert.Equal(t, checkSum+".exe", conf.FileName)
	})

	t.Run("GetExtFromHeader_2", func(t *testing.T) {
		conf := downloader.Config{FileName: "", KeepOriginalName: false}
		r := httptest.NewRecorder()
		resp := r.Result()
		req := httptest.NewRequest("GET", "http://example.com/foo.exe", nil)
		resp.Request = req
		resp.Header.Add("content-disposition", "attachment;")
		msg := []byte(resp.Request.URL.String())
		checkSum := utils.GetChecksum(msg)
		generateFileName(&conf, resp)
		assert.Equal(t, checkSum+".exe", conf.FileName)
	})

	t.Run("GetExtFromHeader_3", func(t *testing.T) {
		conf := downloader.Config{FileName: "", KeepOriginalName: false}
		r := httptest.NewRecorder()
		resp := r.Result()
		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		resp.Request = req
		resp.Header.Add("content-disposition", "attachment;")
		msg := []byte(resp.Request.URL.String())
		checkSum := utils.GetChecksum(msg)
		generateFileName(&conf, resp)
		assert.Equal(t, checkSum, conf.FileName)
	})

	t.Run("KeepOriginalName", func(t *testing.T) {
		conf := downloader.Config{FileName: "", KeepOriginalName: true}
		r := httptest.NewRecorder()
		resp := r.Result()
		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		resp.Request = req
		generateFileName(&conf, resp)
		assert.Equal(t, "foo", conf.FileName)
	})

	t.Run("KeepOriginalNameWithHeader_1", func(t *testing.T) {
		conf := downloader.Config{FileName: "", KeepOriginalName: true}
		r := httptest.NewRecorder()
		resp := r.Result()
		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		resp.Request = req
		resp.Header.Add("content-disposition", "attachment; filename=foo.exe")
		generateFileName(&conf, resp)
		assert.Equal(t, "foo.exe", conf.FileName)
	})

	t.Run("KeepOriginalNameWithHeader_2", func(t *testing.T) {
		conf := downloader.Config{FileName: "", KeepOriginalName: true}
		r := httptest.NewRecorder()
		resp := r.Result()
		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		resp.Request = req
		resp.Header.Add("content-disposition", "attachment;")
		generateFileName(&conf, resp)
		assert.Equal(t, "foo", conf.FileName)
	})

}

func Test_serviceImpl_GetChecksum(t *testing.T) {
	ctrl := gomock.NewController(t)

	clientMock := cMock.NewMockClientService(ctrl)
	clientMock.EXPECT().Do(gomock.Any()).Return(nil, errors.New("Error"))

	resp1 := &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader("Read"))}
	clientMock1 := cMock.NewMockClientService(ctrl)
	clientMock1.EXPECT().Do(gomock.Any()).Return(resp1, nil)

	type fields struct {
		client webClient.ClientService
		log    logging.Logger
	}
	type args struct {
		conf *downloader.Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    io.Reader
		want1   string
		wantErr bool
	}{
		{
			name: "TC1",
			fields: fields{
				client: clientMock,
				log:    logging.GetLoggerFactory().Get(),
			},
			args: args{
				conf: &downloader.Config{URL: "http://test", DownloadLocation: downloadLocation, FileName: fileName},
			},
			want:    nil,
			want1:   "",
			wantErr: true,
		},
		{
			name: "TC2",
			fields: fields{
				client: clientMock,
				log:    logging.GetLoggerFactory().Get(),
			},
			args: args{
				conf: &downloader.Config{URL: "http://test", DownloadLocation: downloadLocation, FileName: fileName, CheckSum: "123"},
			},
			want:    nil,
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := serviceImpl{
				client: tt.fields.client,
				log:    tt.fields.log,
			}
			got, got1, err := s.getChecksum(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("serviceImpl.GetChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("serviceImpl.GetChecksum() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("serviceImpl.GetChecksum() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
