package cloud

//Resource represents a cloud resource
type Resource struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	ServiceType string              `json:"service"`
	Location    string              `json:"location"`
	Hierarchies []ResourceHierarchy `json:"hierarchies"`
	AvailabilityState string `json:availabilitystatus`
}

//ResourceHierarchy represents a Resource hierarchy
type ResourceHierarchy struct {
	Title string `json:"title"`
	Name  string `json:"name"`
	ID    string `json:"id"`
	Level int    `json:"level"`
}
