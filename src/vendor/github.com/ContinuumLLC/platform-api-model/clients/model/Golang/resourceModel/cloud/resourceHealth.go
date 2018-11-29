package cloud

//ResourceHealthDetails represents a cloud resource
type ResourceHealthDetails struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	AvailabilityStatus string              `json:"availabilitystatus"`
	ImpactingEvents    []ImpactingEvent    `json:"impactingevents"`
	RecommendedActions []RecommendedAction `json:"recommendedactions"`
}

//ImpactingEvent details of services
type ImpactingEvent struct {
	Summary           string `json:"summary"`
	Description       string `json:"description"`
	Event             string `json:"event"`
	EventType         string `json:"eventtype"`
	Location          string `json:"location"`
	OccurrenceTimeUTC string `json:"occurrencetimeutc"`
	LastUpdateTimeUTC string `json:"lastupdatetimeutc"`
	ResolutionETA     string `json:"resolutioneta"`
	Status            string `json:"status"`
}

//RecommendedAction actions
type RecommendedAction struct {
	Description string `json:"description"`
	URL         string `json:"url"`
	URLTitle    string `json:"urltitle"`
}
