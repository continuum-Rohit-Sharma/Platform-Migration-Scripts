package cloud

//Hierarchy represents an hierarchy of the cloud vendors product suite
type Hierarchy struct {
	Title  string           `json:"title"`
	Level  int              `json:"level"`
	Values []HierarchyValue `json:"values"`
}

//HierarchyValue represents a single value in the list of values for a hierarchy
type HierarchyValue struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}
