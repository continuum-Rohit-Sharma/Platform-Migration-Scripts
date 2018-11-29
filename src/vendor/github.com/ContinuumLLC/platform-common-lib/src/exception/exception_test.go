package exception

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestException(t *testing.T) {
	exc := New("SetErrorCode", nil)

	if exc.GetErrorCode() != "SetErrorCode" {
		t.Error("ErrorCode Match Failed")
	}

	if exc.GetInnerError() != nil {
		t.Error("Inner Error Nil Check Failed")
	}

	if exc.Error() != "SetErrorCode" {
		t.Error("Error() string Match Failed")
	}

	outer := New("OuterErrorCode", exc)

	if outer.GetErrorCode() != "OuterErrorCode" {
		t.Error("Outer Error Code Match")
	}

	if outer.GetInnerError() != exc {
		t.Error("Inner Exception Match Failed")
	}

	if outer.Error() != "OuterErrorCode" {
		t.Error("Outer Error() string Match Failed: " + exc.Error())
	}
}

func TestFormatVPlus(t *testing.T) {
	str, exc, _ := testFormatCreateError("%+v", false)

	testFormatContains(t, str, exc.errorCode, true)
	testFormatContains(t, str, "exception.testFormatCreateError", true)
	testFormatContains(t, str, "exception_test.go", true)
	testFormatContains(t, str, "testing.tRunner", true)
	testFormatContains(t, str, "runtime.goexit", true)
}

func testFormatCreateError(format string, createInner bool) (str string, outer *exceptionImpl, inner error) {
	buf := new(bytes.Buffer)
	if createInner {
		inner = New("Inner", nil)
	}
	err := New("Outer", inner)
	outer = err.(*exceptionImpl)
	fmt.Fprintf(buf, format, outer)
	str = string(buf.Bytes())
	return
}

func testFormatContains(t *testing.T, str string, contains string, result bool) {
	if result != strings.Contains(str, contains) {
		t.Errorf("Failed Contains[%s] : %t", contains, result)
		t.Errorf(str)
	}
}

// This test case also checks for the first stack trace that is printed by the exception package.
// Failure in this test case, specially because of 2nd and 3rd call to "testFormatContains"
// implies that the stackTrace structure has changed.
// Take a look at exception.cIgnoreInitialCallersCount to fix this issue.
// Also, ensure that all "New" methods create the exeptionImpl instance with the same call stack.
// If the call stack leading to "newException" is different for different methods,
// this test case is likely to fail.
func TestFormatV(t *testing.T) {
	str, exc, _ := testFormatCreateError("%v", false)

	testFormatContains(t, str, exc.errorCode, true)
	testFormatContains(t, str, "exception.testFormatCreateError", true)
	testFormatContains(t, str, "exception_test.go", true)
}

func TestFormatS(t *testing.T) {
	str, exc, _ := testFormatCreateError("%s", false)

	testFormatContains(t, str, exc.errorCode, true)
	testFormatContains(t, str, "exception.testFormatCreateError", false)
	testFormatContains(t, str, "exception_test.go", false)
}

func TestFormatInnerV(t *testing.T) {
	str, exc, _ := testFormatCreateError("%v", true)

	testFormatContains(t, str, exc.errorCode, true)
	testFormatContains(t, str, "TestFormatInnerV", false)
	testFormatContains(t, str, "Inner", true)
}

func TestFormatInnerVPlus(t *testing.T) {
	str, exc, _ := testFormatCreateError("%+v", true)

	testFormatContains(t, str, exc.errorCode, true)
	testFormatContains(t, str, "TestFormatInnerVPlus", true)
	testFormatContains(t, str, "Inner", true)
}
