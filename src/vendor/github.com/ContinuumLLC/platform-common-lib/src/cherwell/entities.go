package cherwell

import (
	"fmt"
)

// Config represents configuration data for connection to CSM
type Config struct {
	ClientID string `json:"client_id"` // ClientID is an API client key for the client making the token request.
	Password string `json:"password"`
	UserName string `json:"user_name"`
	AuthMode string `json:"auth_mode"`
	Host     string `json:"host"` // f.e. https://continuumdev.cherwellondemand.com/CherwellAPI
}

type businessObjectResponse struct {
	BusinessObject
	ErrorData
}

// tokenResponse represents response from Cherwell API after sending request: POST/token request
type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"as:client_id"`
	Expires      string `json:".expires"`
	ExpiresIn    int    `json:"expires_in"`
	Issued       string `json:".issued"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Username     string `json:"username"`
}

// errorResponse contains information about errors which can occur in Cherwell API
type errorResponse struct {
	Err              string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e errorResponse) Error() string {
	return fmt.Sprintf("Error: %v\nDescription: %v\n", e.Err, e.ErrorDescription)
}

// FieldTemplateItem represents structure of the field for Business Object
type FieldTemplateItem struct {
	Dirty       bool   `json:"dirty"`
	DisplayName string `json:"displayName"`
	FieldID     string `json:"fieldId"`
	HTML        string `json:"html"`
	Name        string `json:"name"`
	Value       string `json:"value"`
}

// saveUpdateResponse represents structure of response which Cherwell API sent after request from client to save or update
// Business Object
type saveUpdateResponse struct {
	ErrorData
	PublicID              string                 `json:"busObPublicId"`
	RecordID              string                 `json:"busObRecId"`
	CacheKey              string                 `json:"cacheKey"`
	FieldValidationErrors []FieldValidationError `json:"fieldValidationErrors"`
	NotificationTriggers  []NotificationTrigger  `json:"notificationTriggers"`
}

// saveUpdateBatchResponse represents structure of response which Cherwell API sent after request from client to save or
// update an array of Business Objects in a batch.
type saveUpdateBatchResponse struct {
	Responses []saveUpdateResponse `json:"responses"`
	ErrorData
}

// FieldValidationError represents errors (about Business Object's fields) which occurs after sending invalid request
// to Cherwell API
type FieldValidationError struct {
	Error     string `json:"error"`
	ErrorCode string `json:"errorCode"`
	FieldID   string `json:"fieldId"`
}

// NotificationTrigger contains information for triggering events for notification
type NotificationTrigger struct {
	SourceType   string `json:"sourceType"`
	SourceID     string `json:"sourceId"`
	SourceChange string `json:"sourceChange"`
	Key          string `json:"key"`
}

// saveUpdateRequest represents request to Cherwell API for creating and updating of Business Object
type saveUpdateRequest struct {
	BusinessObject
	CacheKey   string `json:"cacheKey"`
	CacheScope string `json:"cacheScope"` // CacheScope - Enum:	"Tenant", "User", "Session"
	Persist    bool   `json:"persist"`
}

// saveUpdateBatchRequest represents request to Cherwell API for creating and updating of Business Objects in a batch
type saveUpdateBatchRequest struct {
	SaveRequests []saveUpdateRequest `json:"saveRequests"`
	StopOnError  bool                `json:"stopOnError"`
}

// BusinessObjectInfo represents request to Cherwell API for reading and deleting of Business Object
type BusinessObjectInfo struct {
	ID       string `json:"busObId"`
	PublicID string `json:"busObPublicId"`
	RecordID string `json:"busObRecId"`
}

// batchReadRequest represents request to Cherwell API for reading and deleting of Business Object
type batchReadRequest struct {
	ReadRequests []BusinessObjectInfo `json:"readRequests"`
	StopOnError  bool                 `json:"stopOnError"`
}

// batchReadResponse represents response sent by Cherwell API about a batch of Business Object records that includes a list
// of field record IDs, display names, and values for each record
type batchReadResponse struct {
	Responses []businessObjectResponse `json:"responses"`
}

// deleteResponse represents response sent by Cherwell API about deleted Business Object
type deleteResponse struct {
	BusinessObjectInfo
	ErrorData
}

// batchDeleteResponse represents response sent by Cherwell API about a batch of Business Objects which were deleted
type batchDeleteResponse struct {
	Responses []deleteResponse `json:"responses"`
}

// batchDeleteRequest represents request to Cherwell API for deleting a batch of Business Objects
type batchDeleteRequest struct {
	DeleteRequests []BusinessObjectInfo `json:"deleteRequests"`
	StopOnError    bool                 `json:"stopOnError"`
}
