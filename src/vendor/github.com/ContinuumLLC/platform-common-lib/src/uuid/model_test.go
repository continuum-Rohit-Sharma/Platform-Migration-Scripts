package uuid

import (
	"fmt"
	"testing"
)

func TestSatoriUUID(t *testing.T) {
	oneUUID, err := NewRandomUUID()
	if err != nil {
		t.Error("Got an Error")
		return
	}
	fmt.Println(oneUUID.String())
}

func TestSatoriParse(t *testing.T) {
	invalidUUID := "abcdefghi"
	_, err := ParseUUID(invalidUUID)
	if err == nil {
		t.Error(err)
	}

	oneUUID, _ := NewRandomUUID()
	_, err = ParseUUID(oneUUID.String())
	if err != nil {
		t.Error(err)
	}
}
