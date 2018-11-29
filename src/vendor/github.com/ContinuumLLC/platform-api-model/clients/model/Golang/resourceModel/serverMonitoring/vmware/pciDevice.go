package vmware

// PCIDevice describes Hardware PCI device of the Host system
type PCIDevice struct {
	Name   string `json:"name,omitempty"`
	Vendor string `json:"vendor,omitempty"`
}
