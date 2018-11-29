package json

import (
	"bytes"
	"strings"
	"testing"

	"fmt"

	exc "github.com/ContinuumLLC/platform-common-lib/src/exception"
)

func TestReadFileBlankFilePath(t *testing.T) {
	var conf Kafkaconfig
	err := FactoryJSONImpl{}.GetDeserializerJSON().ReadFile(&conf, "")
	exce, ok := err.(exc.Exception)
	if ok {
		if exce.GetErrorCode() != ErrJSONEmptyFilePath {
			t.Errorf("Expected JSONEmptyFilePath but got %v", err)
		}
	} else {
		t.Error("Expecting Exception type")
	}
}

func TestReadFileWrongPath(t *testing.T) {
	var conf Kafkaconfig
	err := FactoryJSONImpl{}.GetDeserializerJSON().ReadFile(&conf, "wrapper_test.json")
	exce, ok := err.(exc.Exception)
	if ok {
		if exce.GetErrorCode() != ErrJSONInvalidFilePathOrUnableToRead {
			t.Errorf("Expected JSONInvalidFilePathOrUnableToRead but got %v", exce)
		}
	} else {
		t.Error("Expecting Exception type")
	}
}

func TestReadStringBlank(t *testing.T) {
	var conf Kafkaconfig
	err := deserializerJSONImpl{}.ReadString(&conf, "")
	exce, ok := err.(exc.Exception)
	if ok {
		if exce.GetErrorCode() != ErrJSONBlankString {
			t.Errorf("Expected ErrJSONBlankString but got %v", exce)
		}
	} else {
		t.Error("Expecting Exception type")
	}
}

func TestWriteFile(t *testing.T) {
	type test struct {
	}
	err := serializerJSONImpl{}.WriteFile("test.txt", &test{})
	if err != nil {
		t.Error("Unexpected error")
	}
}

func TestWriteFileEmptyFile(t *testing.T) {
	type test struct {
	}
	err := serializerJSONImpl{}.WriteFile("", &test{})
	if err == nil || !strings.HasPrefix(err.Error(), "JSONEmptyFilePath") {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestWriteByteStream(t *testing.T) {
	type test struct {
		s string
	}
	var a = "initial"
	var abyte = []byte(`initial`)

	bStream, err := serializerJSONImpl{}.WriteByteStream(&a)
	if err != nil {
		t.Error("Unexpected error")
	}
	fmt.Println(bStream)
	fmt.Println(a == string(bStream))
	//TODO comparison should be not equal to
	if bytes.Compare(bStream, abyte) == 0 {
		// a less b
		t.Errorf("Unexpected error ... %v", string(bStream))
	}
	if string(bStream) == "initial" {
		// a less b
		t.Errorf("Unexpected string returned, Expected %s, Returned %s", a, string(bStream))
	}
}
