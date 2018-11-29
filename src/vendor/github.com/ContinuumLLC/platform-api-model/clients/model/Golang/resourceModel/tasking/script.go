package tasking

import "time"

// Script structure contains data about particular script
type Script struct {
	ID                       string    `json:"id"                       valid:"optional,uuid"`
	PartnerID                string    `json:"partnerId"                valid:"unsettableByUsers"`
	Name                     string    `json:"name"                     valid:"requiredForUsers"`
	Description              string    `json:"description"              valid:"requiredForUsers"`
	Content                  string    `json:"content"                  valid:"requiredForUsers,base64"`
	Engine                   string    `json:"engine"                   valid:"requiredForUsers"`
	EngineMaxVersion         int       `json:"engineMaxVersion"         valid:"requiredForUsers"`
	Categories               []string  `json:"categories"               valid:"requiredForUsers"`
	OutputFormatType         string    `json:"outputFormatType"         valid:"optional"`
	SendIdenticalCounter     int       `json:"sendIdenticalCounter"     valid:"optional"`
	SuppressIdenticalCounter int       `json:"suppressIdenticalCounter" valid:"optional"`
	CreatedAt                time.Time `json:"createdAt"                valid:"unsettableByUsers"`
	CreatedBy                string    `json:"createdBy"                valid:"unsettableByUsers"`
	UpdatedAt                time.Time `json:"updatedAt"                valid:"unsettableByUsers"`
	UpdatedBy                string    `json:"updatedBy"                valid:"unsettableByUsers"`
	JSONSchema               string    `json:"JSONSchema"               valid:"optional"`
	Tags                     []string  `json:"tags"                     valid:"requiredForUsers"`
	Sequence                 bool      `json:"sequence"                 valid:"-"`
	ExpectedExecutionTimeSec int       `json:"expectedExecutionTimeSec" valid:"requiredForUsers"`
	SuccessMessage           string    `json:"successMessage"           valid:"requiredForExternal"`
	FailureMessage           string    `json:"failureMessage"           valid:"requiredForExternal"`
	UISchema                 string    `json:"UISchema"                 valid:"optional"`
	Internal                 bool      `json:"internal"                 valid:"-"`
	Deleted                  bool      `json:"-"                        valid:"-"`
	NOCVisibleOnly           bool      `json:"NOCVisibleOnly"           valid:"-"`
}
