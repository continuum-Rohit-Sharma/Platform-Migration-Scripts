package protocol

import "testing"

// TestHeaders ...
func TestHeaders(t *testing.T) {
	headers := make(Headers)
	h := &headers
	hkey1 := HeaderKey("key1")
	hkey2 := HeaderKey("key2")
	h.SetKeyValue(hkey1, "value1")
	checkHeaderValue(t, h, "key1", "value1")
	h.SetKeyValues(hkey2, []string{"one", "two"})
	checkHeaderValue(t, h, "key2", "one")

	gotValues := h.GetKeyValues(hkey2)
	if gotValues[0] != "one" && gotValues[1] != "two" {
		t.Error("Got Incorrect Values")
		t.Error(gotValues)
	}
}

func checkHeaderValue(t *testing.T, h *Headers, key HeaderKey, value string) {
	valArr, ok := (*h)[key]
	if !ok {
		t.Error("key not found")
		return
	}
	if valArr[0] != value || valArr[0] != h.GetKeyValue(key) {
		t.Error("value doesn't match")
		return
	}
}
