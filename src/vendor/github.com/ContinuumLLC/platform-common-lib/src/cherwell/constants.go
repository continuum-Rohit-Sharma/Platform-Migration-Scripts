package cherwell

// Group of constants related to BusinessObject swagger tag
// it provides uri patterns for formatting with required values
const (
	// These are common business objects IDs
	// incidentID              = "6dd53665c0c24cab86870a21cf6434ae"
	// problemID               = "9344be92d5b7b4c290437c4110bc5b7147c9e3e98a"
	// serviceID               = "9366b3bb9e94d86b5d5f434ff3b657c4ec5bfa3bb3"
	// taskID                  = "9355d5ed41e384ff345b014b6cb1c6e748594aea5b"
	// changeRequestID         = "934ec7a1701c451ce57f2c43bfbbe2e46fe4843f81"
	// statusID                = "5eb3234ae1344c64a19819eda437f18d"
	// descriptionID           = "252b836fc72c4149915053ca1131d138"
	// shortDescriptionID      = "93e8ea93ff67fd95118255419690a50ef2d56f910c"
	// priorityID              = "83c36313e97b4e6b9028aff3b401b71c"
	// customerID              = "933bd530833c64efbf66f84114acabb3e90c6d7b8f"
	// callSourceID            = "93670bdf8abe2cd1f92b1f490a90c7b7d684222e13"
	// ownedByID               = "9339fc404e4c93350bf5be446fb13d693b0bb7f219"

	createUpdateBOEndpoint  = "/api/V1/savebusinessobject"
	createUpdateBOsEndpoint = "/api/V1/savebusinessobjectbatch"
	deleteBOByPubIDEndpoint = "/api/V1/deletebusinessobject/busobid/%v/publicid/%v"
	deleteBOByRecIDEndpoint = "/api/V1/deletebusinessobject/busobid/%v/busobrecid/%v"
	deleteBOsEndpoint       = "/api/V1/deletebusinessobjectbatch"
	getBOByPubIDEndpoint    = "/api/V1/getbusinessobject/busobid/%v/publicid/%v"
	getBOByRecIDEndpoint    = "/api/V1/getbusinessobject/busobid/%v/busobrecid/%v"
	getBOsEndpoint          = "/api/V1/getbusinessobjectbatch"
	retryCount              = 3
)

// searchEndpoint is a constant declared with uri for performing search requests
const searchEndpoint = "/api/V1/getsearchresults"

// tokenEndpoint is a constant declared with uri for token obtaining
const tokenEndpoint = "/token"

// passwordGrantType is a grant type value to be sent in token refresh request
const passwordGrantType = "password"

const attachmentUploadPath = "/api/V1/uploadbusinessobjectattachment"
const attachmentGetPath = "/api/V1/getbusinessobjectattachments/busobid/%s/busobrecid/%s/type/%s/attachmenttype/%s"
const attachmentUploadPathTemplate = "/filename/%s/busobid/%s/busobrecid/%s/offset/%d/totalsize/%s"
const attachmentDeletePath = "/api/V1/removebusinessobjectattachment/attachmentid/%s/busobid/%s/busobrecid/%s"
const attachmentDownloadPath = "/api/V1/getbusinessobjectattachment/attachmentid/%s/busobid/%s/busobrecid/%s"
