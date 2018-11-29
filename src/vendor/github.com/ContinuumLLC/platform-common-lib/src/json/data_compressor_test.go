package json

import (
	"reflect"
	"testing"
)

func Test_gzipCompress(t *testing.T) {
	type args struct {
		rawData []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Test_gzipCompress :given nil input data",
			args{nil},
			true,
		},
		{
			"Test_gzipCompress :given empty data",
			args{[]byte("")},
			true,
		},
		{
			"Test_gzipCompress :given valid data",
			args{[]byte(`data:"status"`)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := gzipCompress(tt.args.rawData)
			if (err != nil) != tt.wantErr {
				t.Errorf("gzipCompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCompress(t *testing.T) {
	type args struct {
		compType CompressionType
		rawData  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"TestCompress :given valid data and supported gzip compression type",
			args{GZIP, []byte(`data:"status"`)},
			false,
		},
		{
			"TestCompress :given valid data and notsupported compression type",
			args{"not", []byte(`data:"status"`)},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Compress(tt.args.compType, tt.args.rawData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetCompressionType(t *testing.T) {
	type args struct {
		rawData string
	}
	tests := []struct {
		name    string
		args    args
		want    CompressionType
		wantErr bool
	}{
		{
			"TestGetCompressionType :given valid data as gzip compression type",
			args{"gzip"},
			GZIP,
			false,
		},
		{
			"TestGetCompressionType :given valid data as gzip compression type but with spaces",
			args{"gzip                           "},
			GZIP,
			false,
		},
		{
			"TestGetCompressionType :given invalid data as deflate compression type which is not supported",
			args{"deflate"},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCompressionType(tt.args.rawData)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestGetCompressionType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestGetCompressionType() got = %v, wantErr %v", got, tt.want)
				return
			}
		})
	}
}
