package cherwell

import (
	"fmt"
	"net/http"
)

const (
	// OpEqual is a constant for determining equality operation
	OpEqual = "eq"
)

// PromptValue is a structure representing prompt value for search request
type PromptValue struct {
	ID                       string      `json:"busObId,omitempty"`
	CollectionStoreEntireRow string      `json:"collectionStoreEntireRow,omitempty"`
	CollectionValueField     string      `json:"collectionValueField,omitempty"`
	FieldID                  string      `json:"fieldId,omitempty"`
	ListReturnFieldID        string      `json:"listReturnFieldId,omitempty"`
	PromptID                 string      `json:"promptId,omitempty"`
	Value                    interface{} `json:"value,omitempty"`
	ValueIsRecID             bool        `json:"valueIsRecId,omitempty"`
}

// Sorting is a structure representing sorting for search request
type Sorting struct {
	FieldID       string `json:"fieldId,omitempty"`
	SortDirection int    `json:"sortDirection,omitempty"`
}

// Filter is a structure representing rules for filtering by business object value
type Filter struct {
	FieldID  string `json:"fieldId"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// SearchRequest request is a structure representing request for endpoint "/api/V1/getsearchresults"
type SearchRequest struct {
	Association        string         `json:"association,omitempty"`
	ID                 string         `json:"busObId,omitempty"`
	CustomGridDefID    string         `json:"customGridDefId,omitempty"`
	DateTimeFormatting string         `json:"dateTimeFormatting,omitempty"`
	FieldID            string         `json:"fieldId,omitempty"`
	Scope              string         `json:"scope,omitempty"`
	ScopeOwner         string         `json:"scopeOwner,omitempty"`
	SearchID           string         `json:"searchId,omitempty"`
	SearchName         string         `json:"searchName,omitempty"`
	SearchText         string         `json:"searchText,omitempty"`
	IncludeAllFields   bool           `json:"includeAllFields,omitempty"`
	IncludeSchema      bool           `json:"includeSchema,omitempty"`
	PageNumber         int            `json:"pageNumber,omitempty"`
	PageSize           int            `json:"pageSize,omitempty"`
	Sorting            []*Sorting     `json:"sorting,omitempty"`
	PromptValues       []*PromptValue `json:"promptValues,omitempty"`
	Fields             []string       `json:"fields,omitempty"`
	Filters            []*Filter      `json:"filters,omitempty"`
}

type searchResultsResponse struct {
	Responses []*businessObjectResponse `json:"businessObjects"`
	ErrorData
}

var filterOperators = map[string]struct{}{
	"eq": {},
	"lt": {},
	"gt": {},
	"contains": {},
	"startswith": {},
}

// NewSearchRequest creates new SearchRequest with ID with specific fields in response
func NewSearchRequest(id string) *SearchRequest {
	return &SearchRequest{ID: id, Filters: []*Filter{}, IncludeAllFields: true}
}

// AddFilter adds new filter to SearchRequest
func (sc *SearchRequest) AddFilter(filedID, operator, value string) error{
	_, ok := filterOperators[operator]
	if !ok {
		return InvalidFilterOperator{Message: "invalid filter operator"}
	}
	f := &Filter{
		FieldID:  filedID,
		Operator: operator,
		Value:    value,
	}
	sc.Filters = append(sc.Filters, f)

	return nil
}

// SetSpecificFields sets specific fields to be returned in BO response
func (sc *SearchRequest) SetSpecificFields(requiredFieldIDs []string) {
	if len(requiredFieldIDs) == 0 {
		return
	}
	sc.IncludeAllFields = false
	sc.Fields = requiredFieldIDs
}

// AppendSpecificFields appends specific fields to be returned in BO response
func (sc *SearchRequest) AppendSpecificFields(requiredFieldIDs ...string) {
	if len(requiredFieldIDs) == 0 {
		return
	}
	sc.IncludeAllFields = false
	sc.Fields = append(sc.Fields, requiredFieldIDs...)
}

// FindBoInfos sets SearchRequest for returning only BoInfos in response
func (c *Client) FindBoInfos(req SearchRequest) ([]BusinessObject, error) {
	input := SearchRequest{ID: req.ID, Filters: req.Filters, IncludeAllFields: false, Fields: []string{""}}

	return c.Find(input)
}

// Find returns all business objects defined by SearchRequest
func (c *Client) Find(req SearchRequest) ([]BusinessObject, error) {
	var resp searchResultsResponse
	err := c.performRequest(http.MethodPost, searchEndpoint, &req, &resp)
	if err != nil {
		return nil, fmt.Errorf("find: %s", err)
	}

	if resp.HasError {
		err = resp.GetErrorObject()
		return nil, err
	}

	errs := NewErrorSet()

	var bos []BusinessObject
	for _, b := range resp.Responses {
		if b.HasError {
			errs.Add(&CherwellError{
				Code: b.ErrorCode, Message: b.ErrorMessage,
			})
		}
		bos = append(bos, BusinessObject{
			BusinessObjectInfo: BusinessObjectInfo{ID: b.ID, PublicID: b.PublicID, RecordID: b.RecordID},
			Fields:             b.Fields,
		})
	}

	if !errs.IsEmpty() {
		return nil, errs
	}

	return bos, nil
}
