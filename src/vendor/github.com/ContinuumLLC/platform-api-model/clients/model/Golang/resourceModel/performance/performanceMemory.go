package performance

import "time"

//PerformanceMemory is the struct definition of /resources/performance/performanceMemory
type PerformanceMemory struct {
	CreateTimeUTC                   time.Time `json:"createTimeUTC"`
	CreatedBy                       string    `json:"createdBy"`
	Index                           int       `json:"index"`
	Name                            string    `json:"name"`
	Type                            string    `json:"type"`
	EndpointID                      string    `json:"endpointID"`
	PartnerID                       string    `json:"partnerID"`
	ClientID                        string    `json:"clientID"`
	SiteID                          string    `json:"siteID"`
	PhysicalTotalBytes              int64     `json:"physicalTotalBytes"`
	PhysicalInUseBytes              int64     `json:"physicalInUseBytes"`
	PhysicalAvailableBytes          int64     `json:"physicalAvailableBytes"`
	VirtualTotalBytes               int64     `json:"virtualTotalBytes"`
	VirtualInUseBytes               int64     `json:"virtualInUseBytes"`
	VirtualAvailableBytes           int64     `json:"virtualAvailableBytes"`
	PercentCommittedInUseBytes      int64     `json:"percentCommittedInUseBytes"`
	CommittedBytes                  int64     `json:"committedBytes"`
	FreeSystemPageTableEntriesBytes int64     `json:"freeSystemPageTableEntriesBytes"`
	PoolNonPagedBytes               int64     `json:"poolNonPagedBytes"`
	PagesPerSecondBytes             int64     `json:"pagesPerSecondBytes"`
	SwapinPerSecondBytes            int64     `json:"swapinPerSecondBytes"`
	SwapoutPerSecondBytes           int64     `json:"swapoutPerSecondBytes"`
	PagesOutputPerSec               int64     `json:"pagesOutputPerSec"`
}
