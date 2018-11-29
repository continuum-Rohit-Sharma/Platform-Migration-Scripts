package json

import (
	"os"
	"testing"
)

type test struct {
}

func TestSerialize(t *testing.T) {

	Serialize(os.Stdout, &test{})
}
