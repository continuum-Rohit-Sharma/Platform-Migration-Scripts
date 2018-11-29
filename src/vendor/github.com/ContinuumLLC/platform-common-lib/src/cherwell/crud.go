package cherwell

import (
	"fmt"
	"net/http"
)

// GetByPublicID returns a Business Object record that includes a list of fields and their public IDs, names, and set values.
func (c *Client) GetByPublicID(boID, publicID string) (*BusinessObject, error) {
	resp := new(businessObjectResponse)
	uri := fmt.Sprintf(getBOByPubIDEndpoint, boID, publicID)
	err := c.performRequest(http.MethodGet, uri, nil, resp)
	if err != nil {
		return nil, err
	}

	if resp.HasError {
		return nil, resp.GetErrorObject()
	}

	return &resp.BusinessObject, nil
}

// GetByRecordID returns a Business Object record that includes a list of fields and their record IDs, names, and set values.
// API endpoint: GET /api/V1/getbusinessobject/busobid/{busobid}/busobrecid/{busobrecid}
func (c *Client) GetByRecordID(ID, RecordID string) (*BusinessObject, error) {
	boResp := new(businessObjectResponse)
	path := fmt.Sprintf(getBOByRecIDEndpoint, ID, RecordID)
	err := c.performRequest(http.MethodGet, path, nil, boResp)
	if err != nil {
		return nil, err
	}
	if boResp.HasError {
		return nil, boResp.GetErrorObject()
	}
	return &boResp.BusinessObject, nil
}

// GetBatch returns a batch of Business Object records that includes a list of field record IDs, display names,
// and values for each record. Specify an array of Business Object IDs, record IDs, or public IDs. Use a flag to stop on
// error or continue on error. POST /api/V1/getbusinessobjectbatch
func (c *Client) GetBatch(items []BusinessObjectInfo) ([]BusinessObject, error) {
	bosResp := new(batchReadResponse)
	request := batchReadRequest{ReadRequests: items, StopOnError: false}
	err := c.performRequest(http.MethodPost, getBOsEndpoint, request, bosResp)
	if err != nil {
		return nil, err
	}
	var bos []BusinessObject
	for _, b := range bosResp.Responses {
		if !b.HasError {
			bos = append(bos, b.BusinessObject)
		}
	}
	// checking if any BO has been founded
	if len(bos) == 0 {
		return nil, &RecordNotFound{Message: "Records Not Found"}
	}

	return bos, nil
}

// Save creates or updates an existing Business Object.
// API endpoint: POST /api/V1/savebusinessobject
func (c *Client) Save(bo BusinessObject) (*BusinessObjectInfo, error) {
	saveResp := new(saveUpdateResponse)
	request := saveUpdateRequest{
		BusinessObject: bo,
		CacheScope:     "Session",
		Persist:        true,
	}
	err := c.performRequest(http.MethodPost, createUpdateBOEndpoint, request, saveResp)
	if err != nil {
		return nil, err
	}
	if saveResp.HasError {
		return nil, saveResp.GetErrorObject()
	}
	boInfo := &BusinessObjectInfo{PublicID: saveResp.PublicID, RecordID: saveResp.RecordID}
	return boInfo, nil
}

// SaveBatch creates or updates an array of Business Objects in a batch.
// API endpoint: POST /api/V1/savebusinessobjectbatch
func (c *Client) SaveBatch(bos []BusinessObject) ([]BusinessObjectInfo, error) {
	saveResp := new(saveUpdateBatchResponse)
	var reqs []saveUpdateRequest
	var boInfos []BusinessObjectInfo
	for _, bo := range bos {
		req := saveUpdateRequest{
			BusinessObject: bo,
			Persist:        true,
		}
		reqs = append(reqs, req)
	}
	request := saveUpdateBatchRequest{SaveRequests: reqs, StopOnError: false}
	err := c.performRequest(http.MethodPost, createUpdateBOsEndpoint, request, saveResp)
	if err != nil {
		return nil, err
	}

	for _, resp := range saveResp.Responses {
		if !resp.HasError {
			boInfo := BusinessObjectInfo{PublicID: resp.PublicID, RecordID: resp.RecordID}
			boInfos = append(boInfos, boInfo)
		}
	}
	if len(boInfos) == 0 {
		return nil, &BusObNotValid{Message: "one or more Business Objects have not enough data to store"}
	}

	return boInfos, nil
}

// DeleteByPublicID deletes a single Business Object by it's own public ID.
func (c *Client) DeleteByPublicID(boID, publicID string) (*BusinessObjectInfo, error) {
	delResp := new(deleteResponse)
	path := fmt.Sprintf(deleteBOByPubIDEndpoint, boID, publicID)
	err := c.performRequest(http.MethodDelete, path, nil, delResp)
	if err != nil {
		return nil, err
	}
	if delResp.HasError {
		return nil, delResp.GetErrorObject()
	}
	return &delResp.BusinessObjectInfo, nil
}

// DeleteByRecordID deletes a single Business Object by it's record Id.
// API endpoint: DELETE /api/V1/deletebusinessobject/busobid/{busobid}/busobrecid/{busobrecid}
func (c *Client) DeleteByRecordID(ID, RecordID string) (*BusinessObjectInfo, error) {
	delResp := new(deleteResponse)
	path := fmt.Sprintf(deleteBOByRecIDEndpoint, ID, RecordID)
	err := c.performRequest(http.MethodDelete, path, nil, delResp)
	if err != nil {
		return nil, err
	}
	if delResp.HasError {
		return nil, delResp.GetErrorObject()
	}
	return &delResp.BusinessObjectInfo, nil
}

// DeleteBatch deletes a batch of Business Objects. Specify an array of Business Object IDs, record IDs, or
// public IDs. Use a flag to stop on error or continue on error. POST /api/V1/deletebusinessobjectbatch
func (c *Client) DeleteBatch(itemsToDelete []BusinessObjectInfo) ([]BusinessObjectInfo, error) {
	var boInfos []BusinessObjectInfo
	delResp := new(batchDeleteResponse)
	request := batchDeleteRequest{DeleteRequests: itemsToDelete, StopOnError: false}
	err := c.performRequest(http.MethodPost, deleteBOsEndpoint, request, delResp)
	if err != nil {
		return nil, err
	}
	for _, b := range delResp.Responses {
		if !b.HasError {
			boInfo := b.BusinessObjectInfo
			boInfos = append(boInfos, boInfo)
		}
	}
	// checking if any BO has been founded
	if len(boInfos) == 0 {
		return nil, &RecordNotFound{Message: "Records Not Found"}
	}
	return boInfos, nil
}
