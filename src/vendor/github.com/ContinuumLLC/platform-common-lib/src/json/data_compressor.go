package json

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"strings"
)

/**CompressionType holds the value of supported std. compression techniques */
type CompressionType string

const (
	/**GZIP CompressionType*/
	GZIP CompressionType = "GZIP"
)

//Compress accepts compression type and data to be compressed as input
//returns either- error incase compressionType is empty or not supported -or compressed data as per underlying compression
//technique
func Compress(compType CompressionType, rawData []byte) ([]byte, error) {
	switch compType {
	case GZIP:
		return gzipCompress(rawData)
	default:
		errMsg := fmt.Sprintf("Compression Type : %v is not yet supported.", compType)
		return nil, errors.New(errMsg)
	}
}

//GetCompressionType will return pre defined matched compression type against raw string value provided
//GetCompressionType will return pre defined matched compression type against raw string value provided
func GetCompressionType(rawValue string) (CompressionType, error) {
	val := strings.TrimSpace(rawValue)
	if strings.ToUpper(val) == string(GZIP) {
		return GZIP, nil
	}
	return "", fmt.Errorf("Compression Type : %v is not yet supported.", val)
}

/**
* gzipCompress will use gzip mechanism as compression technique
* returns either- error as per gzip compression impl -or compressed byte array in case of successful compression
 */
func gzipCompress(rawData []byte) ([]byte, error) {
	if len(rawData) <= 0 {
		return nil, errors.New("gzip compress data gets invalid input data")
	}
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	defer gz.Close()
	_, err := gz.Write(rawData)
	if err != nil {
		return nil, err
	}
	if err = gz.Flush(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
