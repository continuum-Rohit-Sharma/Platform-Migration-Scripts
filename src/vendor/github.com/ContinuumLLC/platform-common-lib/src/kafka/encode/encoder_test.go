package encode

import (
	"testing"

	"github.com/Shopify/sarama"
)

func TestGetBytesEncoder(t *testing.T) {
	got := GetBytesEncoder(nil)
	_, ok := got.(sarama.ByteEncoder)
	if !ok {
		t.Error("Invalid sarama.ByteEncoder")
	}

}

func TestGetStringEncoder(t *testing.T) {
	got := GetStringEncoder("")
	_, ok := got.(sarama.StringEncoder)
	if !ok {
		t.Error("Invalid sarama.StringEncoder")
	}
}
