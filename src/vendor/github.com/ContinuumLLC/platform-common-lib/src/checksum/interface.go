package checksum

import (
	"errors"
	"io"
)

//Service interface for methods for checksums
type Service interface {
	Calculate(reader io.Reader) (string, error)
	Validate(reader io.Reader, checksum string) (bool, error)
}

//Type is the checksum type
type Type int

const (
	//NONE is a type used for no check sum calculation and validation
	NONE Type = -1
	//MD5 is a type used for md5 check sum calculation and validation
	MD5 Type = -2
	//SHA1 is a type used for SHA1 check sum calculation and validation
	SHA1 Type = -3
	//ErrChecksumInvalid is returned if validation failed du to invalid value
	ErrChecksumInvalid string = "ErrChecksumInvalid"
	//ErrUnsupportedType is returned if invalid type is supplied
	ErrUnsupportedType string = "ErrUnsupportedChecksumType"
)

//GetService is the method to get the checksum type sum method
func GetService(cType Type) (Service, error) {
	switch cType {
	case MD5:
		return md5Impl{}, nil
	case SHA1:
		return sha1Impl{}, nil
	default:
		return nil, errors.New(ErrUnsupportedType)
	}
}

//GetType returns the int value for Type
func GetType(sType string) Type {
	switch sType {
	case "MD5":
		return MD5
	case "SHA1":
		return SHA1
	default:
		return NONE
	}
}
