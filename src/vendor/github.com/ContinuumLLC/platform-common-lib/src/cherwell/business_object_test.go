package cherwell

import "testing"

func TestAddField(t *testing.T) {
	bo := NewBusinessObject("test")

	bo.AddField(&FieldTemplateItem{
		Dirty:   false,
		FieldID: "1111",
		Name:    "foobar",
		Value:   "test",
	}).MarkFieldsAsDirty()

	field, found := bo.FieldByName("foobar")

	if !found {
		t.Fatal("added field was not returned")
	}

	if !field.Dirty {
		t.Fatal("field should be marked as dirty")
	}

	bo.SetField(&FieldTemplateItem{
		Dirty:   true,
		FieldID: "1111",
		Name:    "foobar",
		Value:   "changed",
	})

	field, found = bo.FieldByName("foobar")
	if !found {
		t.Fatal("field not available after update")
	}

	if field.Value != "changed" {
		t.Fatal("field was not updated after SetField operation")
	}

}
