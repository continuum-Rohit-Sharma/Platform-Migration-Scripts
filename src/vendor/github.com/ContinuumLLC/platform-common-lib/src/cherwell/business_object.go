package cherwell

// BusinessObject is a structure which represents all Business Object in Cherwell API
type BusinessObject struct {
	BusinessObjectInfo
	Fields []FieldTemplateItem `json:"fields"`
}

// FieldByName gets business object by field name.
func (bo *BusinessObject) FieldByName(name string) (*FieldTemplateItem, bool) {
	for index, field := range bo.Fields {
		if field.Name == name {
			return &bo.Fields[index], true
		}
	}

	return nil, false
}

// SetField updates existing field in BO or adds a new one if field doesn't exists
func (bo *BusinessObject) SetField(newField *FieldTemplateItem) *BusinessObject {
	for i, field := range bo.Fields {
		if field.FieldID == newField.FieldID {
			bo.Fields[i] = *newField
			return bo
		}
	}

	return bo.AddField(newField)
}

// AddField adds a new field to business object
func (bo *BusinessObject) AddField(field *FieldTemplateItem) *BusinessObject {
	bo.Fields = append(bo.Fields, *field)
	return bo
}

// MarkFieldsAsDirty marks all fields as dirty
func (bo *BusinessObject) MarkFieldsAsDirty() *BusinessObject {
	for i := range bo.Fields {
		bo.Fields[i].Dirty = true
	}

	return bo
}

// NewBusinessObject returns a new BO instance
func NewBusinessObject(id string) *BusinessObject {
	return &BusinessObject{
		BusinessObjectInfo: BusinessObjectInfo{
			ID: id,
		},
	}
}
