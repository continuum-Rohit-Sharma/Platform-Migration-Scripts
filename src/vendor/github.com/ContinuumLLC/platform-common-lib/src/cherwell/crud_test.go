package cherwell

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	GetByPubIDTestcases = map[string]struct {
		busobID, publicID string
		expected          *BusinessObject
		err               error
		path, method      string
		resp              string
		statusCode        int
	}{
		"Test GetBusinessObjectByPubID success": {
			busobID:  "id_1",
			publicID: "pub_id_1",
			expected: &BusinessObject{
				BusinessObjectInfo: BusinessObjectInfo{
					ID:       "id_1",
					PublicID: "pub_id_1",
					RecordID: "rec_id_1",
				},
				Fields: []FieldTemplateItem{
					{
						Dirty:       true,
						DisplayName: "display_name_1",
						FieldID:     "field_id_1",
						HTML:        `<div id='somediv1'></div>`,
						Name:        "name_1",
						Value:       "value_1",
					},
					{
						Dirty:       false,
						DisplayName: "display_name_2",
						FieldID:     "field_id_2",
						HTML:        `<span id='somediv2'></span>`,
						Name:        "name_2",
						Value:       "value_2",
					},
				},
			},
			err:    nil,
			path:   "/api/V1/getbusinessobject/busobid/id_1/publicid/pub_id_1",
			method: http.MethodGet,
			resp: `{
				"busObId": "id_1",
				"busObPublicId": "pub_id_1",
				"busObRecId": "rec_id_1",
				"fields": [
				  {
					"dirty": true,
					"displayName": "display_name_1",
					"fieldId": "field_id_1",
					"html": "<div id='somediv1'></div>",
					"name": "name_1",
					"value": "value_1"
				  },
				  {
					"dirty": false,
					"displayName": "display_name_2",
					"fieldId": "field_id_2",
					"html": "<span id='somediv2'></span>",
					"name": "name_2",
					"value": "value_2"
				  }
				],
				"errorCode": "",
				"errorMessage": "",
				"hasError": false
			  }`,
			statusCode: http.StatusOK,
		},
		"Test GetBusinessObjectByPubID response with cherwell error": {
			busobID:  "id_wrong",
			publicID: "pub_id_wrong",
			expected: nil,
			err: &CherwellError{
				Code:    GeneralFailureError,
				Message: "Business object with 12 not found Parameter name: BusObId",
			},
			path:   "/api/V1/getbusinessobject/busobid/id_wrong/publicid/pub_id_wrong",
			method: http.MethodGet,
			resp: `{
					"busObId": null,
					"busObPublicId": null,
					"busObRecId": null,
					"fields": [],
					"errorCode": "GENERALFAILURE",
					"errorMessage": "Business object with 12 not found Parameter name: BusObId",
					"hasError": true
				}`,
			statusCode: http.StatusInternalServerError,
		},
		"Test GetBusinessObjectByPubID response with client error": {
			busobID:    "id_wrong_json",
			publicID:   "pub_id_wrong_json",
			expected:   nil,
			err:        errors.New("non-JSON response received, unexpected end of JSON input: {"),
			path:       "/api/V1/getbusinessobject/busobid/id_wrong_json/publicid/pub_id_wrong_json",
			method:     http.MethodGet,
			resp:       `{`,
			statusCode: http.StatusInternalServerError,
		},
	}
	GetByRecordIDTestcases = map[string]struct {
		busobID, recordID string
		expected          *BusinessObject
		err               error
		path, method      string
		resp              string
		statusCode        int
	}{
		"Test GetByPublicID success": {
			busobID:  "id_1",
			recordID: "pub_id_1",
			expected: &BusinessObject{
				BusinessObjectInfo: BusinessObjectInfo{ID: "id_1", PublicID: "pub_id_1", RecordID: "rec_id_1"},
				Fields: []FieldTemplateItem{
					{
						Dirty:       true,
						DisplayName: "display_name_1",
						FieldID:     "field_id_1",
						HTML:        `<div id='somediv1'></div>`,
						Name:        "name_1",
						Value:       "value_1",
					},
					{
						Dirty:       false,
						DisplayName: "display_name_2",
						FieldID:     "field_id_2",
						HTML:        `<span id='somediv2'></span>`,
						Name:        "name_2",
						Value:       "value_2",
					},
				},
			},
			err:    nil,
			path:   "/api/V1/getbusinessobject/busobid/id_1/busobrecid/pub_id_1",
			method: http.MethodGet,
			resp: `{
				"busObId": "id_1",
				"busObPublicId": "pub_id_1",
				"busObRecId": "rec_id_1",
				"fields": [
				  {
					"dirty": true,
					"displayName": "display_name_1",
					"fieldId": "field_id_1",
					"html": "<div id='somediv1'></div>",
					"name": "name_1",
					"value": "value_1"
				  },
				  {
					"dirty": false,
					"displayName": "display_name_2",
					"fieldId": "field_id_2",
					"html": "<span id='somediv2'></span>",
					"name": "name_2",
					"value": "value_2"
				  }
				],
				"links": [
				  {
					"name": "string",
					"url": "string"
				  }
				],
				"errorCode": "",
				"errorMessage": "",
				"hasError": false
			  }`,
			statusCode: http.StatusOK,
		},
		"Test GetByRecordID response with cherwell error": {
			busobID:  "id_wrong",
			recordID: "rec_id_wrong",
			expected: nil,
			err: &RecordNotFound{
				Message: "Business object with 12 not found Parameter name: BusRecId ",
			},
			path:   "/api/V1/getbusinessobject/busobid/id_wrong/busobrecid/rec_id_wrong",
			method: http.MethodGet,
			resp: `{
					"busObId": null,
					"busObPublicId": null,
					"busObRecId": null,
					"fields": [],
					"links": [],
					"errorCode": "GENERALFAILURE",
					"errorMessage": "Business object with 12 not found Parameter name: BusRecId",
					"hasError": true
				}`,
			statusCode: http.StatusInternalServerError,
		},
		"Test GetByRecordID response with client error": {
			busobID:    "id_wrong_json",
			recordID:   "rec_id_wrong_json",
			expected:   nil,
			err:        errors.New("non-JSON response received"),
			path:       "/api/V1/getbusinessobject/busobid/id_wrong_json/busobrecid/rec_id_wrong_json",
			method:     http.MethodGet,
			resp:       `{`,
			statusCode: http.StatusInternalServerError,
		},
	}
	GetBatchTestcases = map[string]struct {
		items        []BusinessObjectInfo
		expected     []BusinessObject
		err          error
		path, method string
		req, resp    string
		statusCode   int
	}{
		"GetBatch success": {
			items: []BusinessObjectInfo{
				{
					ID:       "id_1",
					PublicID: "pub_id_1",
					RecordID: "rec_id_1",
				},
				{
					ID:       "id_2",
					PublicID: "pub_id_2",
					RecordID: "rec_id_2",
				},
				{
					ID:       "id_3",
					PublicID: "pub_id_3",
					RecordID: "rec_id_3",
				},
			},
			expected: []BusinessObject{
				{
					BusinessObjectInfo: BusinessObjectInfo{ID: "id_1", PublicID: "pub_id_1", RecordID: "rec_id_1"},
					Fields: []FieldTemplateItem{
						{
							Dirty:       true,
							DisplayName: "display_name_1",
							FieldID:     "field_id_1",
							HTML:        `<div id='somediv1'></div>`,
							Name:        "name_1",
							Value:       "value_1",
						},
						{
							Dirty:       false,
							DisplayName: "display_name_2",
							FieldID:     "field_id_2",
							HTML:        `<span id='somediv2'></span>`,
							Name:        "name_2",
							Value:       "value_2",
						},
					},
				},
				{
					BusinessObjectInfo: BusinessObjectInfo{ID: "id_2", PublicID: "pub_id_2", RecordID: "rec_id_2"},
					Fields: []FieldTemplateItem{
						{
							Dirty:       false,
							DisplayName: "display_name_2",
							FieldID:     "field_id_2",
							HTML:        `<span id='somediv2'></span>`,
							Name:        "name_2",
							Value:       "value_2",
						},
					},
				},
				{
					BusinessObjectInfo: BusinessObjectInfo{ID: "id_3", PublicID: "pub_id_3", RecordID: "rec_id_3"},
					Fields: []FieldTemplateItem{
						{
							Dirty:       true,
							DisplayName: "display_name_1",
							FieldID:     "field_id_1",
							HTML:        `<div id='somediv1'></div>`,
							Name:        "name_1",
							Value:       "value_1",
						},
					},
				},
			},
			err:    nil,
			path:   "/api/V1/getbusinessobjectbatch",
			method: http.MethodPost,
			resp: `{
				"responses": [
					{
						"busObId": "id_1",
						"busObPublicId": "pub_id_1",
						"busObRecId": "rec_id_1",
						"fields": [
							{
								"dirty": true,
								"displayName": "display_name_1",
								"fieldId": "field_id_1",
								"html": "<div id='somediv1'></div>",
								"name": "name_1",
								"value": "value_1"
							},
							{
								"dirty": false,
								"displayName": "display_name_2",
								"fieldId": "field_id_2",
								"html": "<span id='somediv2'></span>",
								"name": "name_2",
								"value": "value_2"
							}
						],
						"errorCode": "",
						"errorMessage": "",
						"hasError": false
					},
					{
						"busObId": "id_2",
						"busObPublicId": "pub_id_2",
						"busObRecId": "rec_id_2",
						"fields": [
							{
								"dirty": false,
								"displayName": "display_name_2",
								"fieldId": "field_id_2",
								"html": "<span id='somediv2'></span>",
								"name": "name_2",
								"value": "value_2"
							}
						],
						"errorCode": "",
						"errorMessage": "",
						"hasError": false
					},
					{
						"busObId": "id_3",
						"busObPublicId": "pub_id_3",
						"busObRecId": "rec_id_3",
						"fields": [
							{
								"dirty": true,
								"displayName": "display_name_1",
								"fieldId": "field_id_1",
								"html": "<div id='somediv1'></div>",
								"name": "name_1",
								"value": "value_1"
							}
						],
						"errorCode": "",
						"errorMessage": "",
						"hasError": false
					}
				]
			}`,
			statusCode: http.StatusOK,
		},
		"GetBatch one of BOs is invalid": {
			items: []BusinessObjectInfo{
				{
					ID:       "id_1",
					PublicID: "pub_id_1",
					RecordID: "rec_id_1",
				},
				{
					ID:       "id_2",
					PublicID: "pub_id_2",
					RecordID: "rec_id_2",
				},
				{
					ID:       "id_invalid",
					PublicID: "pub_id_invalid",
					RecordID: "rec_id_invalid",
				},
			},
			expected: []BusinessObject{
				{
					BusinessObjectInfo: BusinessObjectInfo{ID: "id_1", PublicID: "pub_id_1", RecordID: "rec_id_1"},
					Fields: []FieldTemplateItem{
						{
							Dirty:       true,
							DisplayName: "display_name_1",
							FieldID:     "field_id_1",
							HTML:        `<div id='somediv1'></div>`,
							Name:        "name_1",
							Value:       "value_1",
						},
						{
							Dirty:       false,
							DisplayName: "display_name_2",
							FieldID:     "field_id_2",
							HTML:        `<span id='somediv2'></span>`,
							Name:        "name_2",
							Value:       "value_2",
						},
					},
				},
				{
					BusinessObjectInfo: BusinessObjectInfo{ID: "id_2", PublicID: "pub_id_2", RecordID: "rec_id_2"},
					Fields: []FieldTemplateItem{
						{
							Dirty:       false,
							DisplayName: "display_name_2",
							FieldID:     "field_id_2",
							HTML:        `<span id='somediv2'></span>`,
							Name:        "name_2",
							Value:       "value_2",
						},
					},
				},
				{
					Fields: []FieldTemplateItem{},
				},
			},
			err: &CherwellError{
				Code:    "GENERALFAILURE",
				Message: "Business object with id_invalid not found\r\nParameter name: BusObId",
			},
			path:   "/api/V1/getbusinessobjectbatch",
			method: http.MethodPost,
			resp: `{
				"responses": [
					{
						"busObId": "id_1",
						"busObPublicId": "pub_id_1",
						"busObRecId": "rec_id_1",
						"fields": [
							{
								"dirty": true,
								"displayName": "display_name_1",
								"fieldId": "field_id_1",
								"html": "<div id='somediv1'></div>",
								"name": "name_1",
								"value": "value_1"
							},
							{
								"dirty": false,
								"displayName": "display_name_2",
								"fieldId": "field_id_2",
								"html": "<span id='somediv2'></span>",
								"name": "name_2",
								"value": "value_2"
							}
						],
						"errorCode": "",
						"errorMessage": "",
						"hasError": false
					},
					{
						"busObId": "id_2",
						"busObPublicId": "pub_id_2",
						"busObRecId": "rec_id_2",
						"fields": [
							{
								"dirty": false,
								"displayName": "display_name_2",
								"fieldId": "field_id_2",
								"html": "<span id='somediv2'></span>",
								"name": "name_2",
								"value": "value_2"
							}
						],
						"errorCode": "",
						"errorMessage": "",
						"hasError": false
					},
					{
						"busObId": null,
						"busObPublicId": null,
						"busObRecId": null,
						"fields": [],
						"links": [],
						"errorCode": "GENERALFAILURE",
						"errorMessage": "Business object with id_invalid not found\r\nParameter name: BusObId",
						"hasError": true
					  }
				]
			}`,
			statusCode: http.StatusOK,
		},
		"GetBatch invalid cherwell response": {
			items: []BusinessObjectInfo{
				{
					ID:       "id_1",
					PublicID: "pub_id_1",
					RecordID: "rec_id_1",
				},
			},
			expected:   nil,
			err:        errors.New("non-JSON response received"),
			path:       "/api/V1/getbusinessobjectbatch",
			method:     http.MethodPost,
			resp:       `{`,
			statusCode: http.StatusOK,
		},
	}
	SaveTestcases = map[string]struct {
		bo             BusinessObject
		expectedBoInfo *BusinessObjectInfo
		err            error
		path, method   string
		resp           string
		statusCode     int
	}{
		"Test Save success": {
			expectedBoInfo: &BusinessObjectInfo{PublicID: "pub_id_1", RecordID: "rec_id_1"},
			bo: BusinessObject{
				BusinessObjectInfo: BusinessObjectInfo{ID: "rec_id_1", PublicID: "pub_id_1", RecordID: "rec_id_1"},
				Fields:             []FieldTemplateItem{},
			},
			err:    nil,
			path:   "/api/V1/savebusinessobject",
			method: http.MethodPost,
			resp: `{
				"busObPublicId": "pub_id_1",
				"busObRecId": "rec_id_1",
				"cacheKey": "string",
				"fieldValidationErrors": [],
				"notificationTriggers": [],
				"errorCode": "",
				"errorMessage": "",
				"hasError": false
			  }`,
			statusCode: http.StatusOK,
		},
		"Test Save response with cherwell error": {
			bo: BusinessObject{
				BusinessObjectInfo: BusinessObjectInfo{ID: "id_wrong", PublicID: "pub_id_wrong", RecordID: "rec_id_wrong"},
				Fields:             []FieldTemplateItem{},
			},
			expectedBoInfo: nil,
			err: &CherwellError{
				Code:    "RecordNotFound",
				Message: "Record not found",
			},
			path:   "/api/V1/savebusinessobject",
			method: http.MethodPost,
			resp: `{
					"busObId": null,
					"busObPublicId": null,
					"busObRecId": null,
					"fields": [],
					"links": [],
					"errorCode": "RecordNotFound",
					"errorMessage": "Record not found",
					"hasError": true
				}`,
			statusCode: http.StatusInternalServerError,
		},
		"Test Save response with client error": {
			expectedBoInfo: nil,
			bo: BusinessObject{
				BusinessObjectInfo: BusinessObjectInfo{ID: "id_wrong_json", PublicID: "pub_id_wrong_json", RecordID: "rec_id_wrong_json"},
				Fields:             []FieldTemplateItem{},
			},
			err:        errors.New("non-JSON response received"),
			path:       "/api/V1/savebusinessobject",
			method:     http.MethodPost,
			resp:       `{`,
			statusCode: http.StatusInternalServerError,
		},
	}
	SaveBatchTestcases = map[string]struct {
		bos             []BusinessObject
		expectedBoInfos []BusinessObjectInfo
		err             error
		path, method    string
		req, resp       string
		statusCode      int
	}{
		"SaveBatch success": {
			expectedBoInfos: []BusinessObjectInfo{
				{PublicID: "pub_id_1", RecordID: "rec_id_1"},
				{PublicID: "pub_id_2", RecordID: "rec_id_2"},
				{PublicID: "pub_id_3", RecordID: "rec_id_3"}},
			bos: []BusinessObject{
				{
					BusinessObjectInfo: BusinessObjectInfo{ID: "id_1", PublicID: "pub_id_1", RecordID: "rec_id_1"},
					Fields:             []FieldTemplateItem{},
				},
				{
					BusinessObjectInfo: BusinessObjectInfo{ID: "id_2", PublicID: "pub_id_2", RecordID: "rec_id_2"},
					Fields:             []FieldTemplateItem{},
				},
				{
					BusinessObjectInfo: BusinessObjectInfo{ID: "id_3", PublicID: "pub_id_3", RecordID: "rec_id_3"},
					Fields:             []FieldTemplateItem{},
				},
			},
			err:    nil,
			path:   "/api/V1/savebusinessobjectbatch",
			method: http.MethodPost,
			resp: `{
				"responses": [
				  {
					"busObPublicId": "pub_id_1",
					"busObRecId": "rec_id_1",
					"cacheKey": "string",
					"fieldValidationErrors": [],
					"notificationTriggers": [],
					"errorCode": "",
					"errorMessage": "",
					"hasError": false
				  },
				  {
					"busObPublicId": "pub_id_2",
					"busObRecId": "rec_id_2",
					"cacheKey": "string",
					"fieldValidationErrors": [],
					"notificationTriggers": [],
					"errorCode": "",
					"errorMessage": "",
					"hasError": false
				  },
				  {
					"busObPublicId": "pub_id_3",
					"busObRecId": "rec_id_3",
					"cacheKey": "string",
					"fieldValidationErrors": [],
					"notificationTriggers": [],
					"errorCode": "",
					"errorMessage": "",
					"hasError": false
				  }
				],
				"errorCode": "",
				"errorMessage": "",
				"hasError": false
			  }`,
			statusCode: http.StatusOK,
		},
		"SaveBatch one of BOs is invalid": {
			expectedBoInfos: []BusinessObjectInfo{{ID: "", PublicID: "pub_id_1", RecordID: "rec_id_1"},
				{ID: "", PublicID: "pub_id_2", RecordID: "rec_id_2"}},
			bos: []BusinessObject{
				{
					BusinessObjectInfo: BusinessObjectInfo{ID: "id_1", PublicID: "pub_id_1", RecordID: "rec_id_1"},
					Fields:             []FieldTemplateItem{},
				},
				{
					BusinessObjectInfo: BusinessObjectInfo{ID: "id_2", PublicID: "pub_id_2", RecordID: "rec_id_2"},
					Fields:             []FieldTemplateItem{},
				},
				{
					Fields: []FieldTemplateItem{},
				},
			},
			err: &CherwellError{
				Code:    "RecordNotFound",
				Message: "Record not found",
			},
			path:   "/api/V1/savebusinessobjectbatch",
			method: http.MethodPost,
			resp: `{
				"responses": [
					{
						"busObPublicId": "pub_id_1",
						"busObRecId": "rec_id_1",
						"cacheKey": "string",
						"fieldValidationErrors": [],
						"notificationTriggers": [],
						"errorCode": "",
						"errorMessage": "",
						"hasError": false
					},
					{
						"busObPublicId": "pub_id_2",
						"busObRecId": "rec_id_2",
						"cacheKey": "string",
						"fieldValidationErrors": [],
						"notificationTriggers": [],
						"errorCode": "",
						"errorMessage": "",
						"hasError": false
					},
					{
						"busObPublicId": null,
						"busObRecId": null,
						"cacheKey": null,
						"fieldValidationErrors": [],
						"notificationTriggers": [],
						"errorCode": "RecordNotFound",
						"errorMessage": "Record not found",
						"hasError": true
					  }
				],
				"errorCode": "RecordNotFound",
				"errorMessage": "Record not found",
				"hasError": true
			}`,
			statusCode: http.StatusOK,
		},
		"GetBatch invalid cherwell response": {
			expectedBoInfos: nil,
			err:             errors.New("non-JSON response received"),
			path:            "/api/V1/savebusinessobjectbatch",
			method:          http.MethodPost,
			resp:            "{",
			statusCode:      http.StatusOK,
		},
	}
	DeleteByPubIDTestcases = map[string]struct {
		busobID, publicID            string
		expectedPubID, expectedRecID string
		err                          error
		path, method                 string
		resp                         string
		statusCode                   int
	}{
		"Test DeleteBusinessObjectByPubID success": {
			busobID:       "id_1",
			publicID:      "pub_id_1",
			expectedPubID: "pub_id_1",
			expectedRecID: "rec_id_1",
			err:           nil,
			path:          "/api/V1/deletebusinessobject/busobid/id_1/publicid/pub_id_1",
			method:        http.MethodDelete,
			resp: `{
				"busObId": "id_1",
				"busObPublicId": "pub_id_1",
				"busObRecId": "rec_id_1",
				"errorCode": "",
				"errorMessage": "",
				"hasError": false
			  }`,
			statusCode: http.StatusOK,
		},
		"Test GetBusinessObjectByPubID response with cherwell error": {
			busobID:  "id_wrong",
			publicID: "pub_id_wrong",
			err: &CherwellError{
				Code:    "GENERALFAILURE",
				Message: "Business object with 12 not found Parameter name: BusObId",
			},
			path:   "/api/V1/deletebusinessobject/busobid/id_wrong/publicid/pub_id_wrong",
			method: http.MethodDelete,
			resp: `{
					"busObId": null,
					"busObPublicId": null,
					"busObRecId": null,
					"errorCode": "GENERALFAILURE",
					"errorMessage": "Business object with 12 not found Parameter name: BusObId",
					"hasError": true
				}`,
			statusCode: http.StatusInternalServerError,
		},
		"Test DeleteBusinessObjectByPubID response with client error": {
			busobID:    "id_wrong_json",
			publicID:   "pub_id_wrong_json",
			err:        errors.New("non-JSON response received, unexpected end of JSON input: "),
			path:       "/api/V1/deletebusinessobject/busobid/id_wrong_json/publicid/pub_id_wrong_json",
			method:     http.MethodGet,
			resp:       `{`,
			statusCode: http.StatusInternalServerError,
		},
	}
	DeleteByRecordIDTestcases = map[string]struct {
		busobID, recordID string
		expectedBoInfo    *BusinessObjectInfo
		err               error
		path, method      string
		resp              string
		statusCode        int
	}{
		"Test DeleteByRecordID success": {
			recordID:       "rec_id_1",
			busobID:        "id_1",
			expectedBoInfo: &BusinessObjectInfo{ID: "id_1", PublicID: "pub_id_1", RecordID: "rec_id_1"},
			err:            nil,
			path:           "/api/V1/deletebusinessobject/busobid/id_1/busobrecid/rec_id_1",
			method:         http.MethodDelete,
			resp: `{
				"busObId": "id_1",
				"busObPublicId": "pub_id_1",
				"busObRecId": "rec_id_1",
				"errorCode": "",
				"errorMessage": "",
				"hasError": false
			  }`,
			statusCode: http.StatusOK,
		},
		"Test DeleteByRecordID response with cherwell error": {
			recordID: "rec_id_wrong",
			busobID:  "id_wrong",
			err: &CherwellError{
				Code:    "GENERALFAILURE",
				Message: "Business object with 12 not found Parameter name: BusObId",
			},
			path:   "/api/V1/deletebusinessobject/busobid/id_wrong/busobrecid/rec_id_wrong",
			method: http.MethodDelete,
			resp: `{
					"busObId": null,
					"busObPublicId": null,
					"busObRecId": null,
					"errorCode": "GENERALFAILURE",
					"errorMessage": "Business object with 12 not found Parameter name: BusObId",
					"hasError": true
				}`,
			statusCode: http.StatusInternalServerError,
		},
		"Test DeleteByRecordID response with client error": {
			recordID:   "rec_id_wrong_json",
			busobID:    "id_wrong_json",
			err:        errors.New("non-JSON response received"),
			path:       "/api/V1/deletebusinessobject/busobid/id_wrong_json/busobrecid/rec_id_wrong_json",
			method:     http.MethodGet,
			resp:       `{`,
			statusCode: http.StatusInternalServerError,
		},
	}
	DeleteBatchTestcases = map[string]struct {
		bos             []BusinessObjectInfo
		expectedBoInfos []BusinessObjectInfo
		err             error
		path, method    string
		resp            string
		statusCode      int
	}{
		"DeleteBatch success": {
			bos: []BusinessObjectInfo{
				{
					ID:       "id_1",
					PublicID: "pub_id_1",
					RecordID: "rec_id_1",
				},
				{
					ID:       "id_2",
					PublicID: "pub_id_2",
					RecordID: "rec_id_2",
				},
				{
					ID:       "id_3",
					PublicID: "pub_id_3",
					RecordID: "rec_id_3",
				},
			},
			expectedBoInfos: []BusinessObjectInfo{
				{PublicID: "pub_id_1", RecordID: "rec_id_1"},
				{PublicID: "pub_id_2", RecordID: "rec_id_2"},
				{PublicID: "pub_id_3", RecordID: "rec_id_3"}},
			err:    nil,
			path:   "/api/V1/deletebusinessobjectbatch",
			method: http.MethodPost,
			resp: `{
				"responses": [
				  {
					"busObPublicId": "pub_id_1",
					"busObRecId": "rec_id_1",
					"errorCode": "",
					"errorMessage": "",
					"hasError": false
				  },
				  {
					"busObPublicId": "pub_id_2",
					"busObRecId": "rec_id_2",
					"errorCode": "",
					"errorMessage": "",
					"hasError": false
				  },
				  {
					"busObPublicId": "pub_id_3",
					"busObRecId": "rec_id_3",
					"errorCode": "",
					"errorMessage": "",
					"hasError": false
				  }
				]
			  }`,
			statusCode: http.StatusOK,
		},
		"DeleteBatch one of BOs is invalid": {
			bos: []BusinessObjectInfo{
				{
					ID:       "id_1",
					PublicID: "pub_id_1",
					RecordID: "rec_id_1",
				},
				{
					ID:       "id_2",
					PublicID: "pub_id_2",
					RecordID: "rec_id_2",
				},
				{
					ID:       "id_invalid",
					PublicID: "pub_id_invalid",
					RecordID: "rec_id_invalid",
				},
			},
			expectedBoInfos: []BusinessObjectInfo{
				{PublicID: "pub_id_1", RecordID: "rec_id_1"},
				{PublicID: "pub_id_2", RecordID: "rec_id_2"},
			},
			err: &CherwellError{
				Code:    "GENERALFAILURE",
				Message: "Business object with id_invalid not found\r\nParameter name: BusObId",
			},
			path:   "/api/V1/deletebusinessobjectbatch",
			method: http.MethodPost,
			resp: `{
				"responses": [
					{
						"busObPublicId": "pub_id_1",
						"busObRecId": "rec_id_1",
						"errorCode": "",
						"errorMessage": "",
						"hasError": false
					},
					{
						"busObPublicId": "pub_id_2",
						"busObRecId": "rec_id_2",
						"errorCode": "",
						"errorMessage": "",
						"hasError": false
					},
					{
						"busObPublicId": null,
						"busObRecId": null,
						"errorCode": "GENERALFAILURE",
						"errorMessage": "Business object with id_invalid not found\r\nParameter name: BusObId",
						"hasError": true
					  }
				]
			}`,
			statusCode: http.StatusOK,
		},
		"DeleteBatch invalid cherwell response": {
			bos: []BusinessObjectInfo{
				{
					ID:       "id_1",
					PublicID: "pub_id_1",
					RecordID: "rec_id_1",
				},
			},
			expectedBoInfos: nil,
			err:             errors.New("non-JSON response received"),
			path:            "/api/V1/deletebusinessobjectbatch",
			method:          http.MethodPost,
			resp:            `{`,
			statusCode:      http.StatusOK,
		},
	}
)

func TestByPubID(t *testing.T) {
	for name, tc := range GetByPubIDTestcases {
		server, mux := newTestServer()
		mux.Handle(tc.path, newMockHandler(tc.method, tc.path, tc.resp, tc.statusCode))
		client, err := NewClient(Config{Host: server.URL}, &http.Client{Transport: &http.Transport{}})
		if err != nil {
			assert.NoError(t, err, "Can not create client: %v", err)
		}

		resp, err := client.GetByPublicID(tc.busobID, tc.publicID)
		assert.Equal(t, tc.expected, resp, "Unexpected result on testcase '%s': ", name)
		if tc.err != nil {
			assert.EqualError(t, err, tc.err.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestGetByRecordID(t *testing.T) {
	for name, tc := range GetByRecordIDTestcases {
		server, mux := newTestServer()
		mux.Handle(tc.path, newMockHandler(tc.method, tc.path, tc.resp, tc.statusCode))
		client, err := NewClient(Config{Host: server.URL}, &http.Client{Transport: &http.Transport{}})
		if err != nil {
			assert.NoError(t, err, "Can not create client: %v", err)
		}
		resp, err := client.GetByRecordID(tc.busobID, tc.recordID)
		assert.Equal(t, tc.expected, resp, "Unexpected result on testcase '%s': ", name)
		if tc.err != nil && err != nil {
			if len(err.Error()) != len(tc.err.Error()) {
				if len(err.Error()) < len(tc.err.Error()) {
					assert.Contains(t, tc.err.Error(), err.Error(), "Actual error: %s does not contain expected error: %s", err.Error(), tc.err.Error())
				} else {
					assert.Contains(t, err.Error(), tc.err.Error(), "Actual error: %s does not contain expected error: %s", err.Error(), tc.err.Error())
				}

			} else {
				assert.EqualError(t, err, tc.err.Error())
			}
		} else {
			assert.NoError(t, err)
		}

		server.CloseClientConnections()
		server.Close()
	}
}

func TestGetBatch(t *testing.T) {
	for _, tc := range GetBatchTestcases {
		server, mux := newTestServer()
		mux.Handle(tc.path, newMockHandler(tc.method, tc.path, tc.resp, tc.statusCode))
		client, err := NewClient(Config{Host: server.URL}, &http.Client{Transport: &http.Transport{}})
		if err != nil {
			assert.NoError(t, err, "Can not create client: %v", err)
		}
		_, err = client.GetBatch(tc.items)
		if tc.err != nil && err != nil {
			if len(err.Error()) != len(tc.err.Error()) {
				assert.Contains(t, err.Error(), tc.err.Error(), "Actual error: %s does not contain expected error: %s", err.Error(), tc.err.Error())
			} else {
				assert.EqualError(t, err, tc.err.Error())
			}
		} else {
			assert.NoError(t, err)
		}

		server.CloseClientConnections()
		server.Close()
	}

}

func TestSave(t *testing.T) {
	for name, tc := range SaveTestcases {
		server, mux := newTestServer()
		mux.Handle(tc.path, newMockHandler(tc.method, tc.path, tc.resp, tc.statusCode))
		client, err := NewClient(Config{Host: server.URL}, &http.Client{Transport: &http.Transport{}})
		if err != nil {
			assert.NoError(t, err, "Can not create client: %v", err)
		}

		actualBoInfo, err := client.Save(tc.bo)
		assert.Equal(t, tc.expectedBoInfo, actualBoInfo, "Unexpected result on testcase '%s': ", name)

		if tc.err != nil && err != nil {
			if len(err.Error()) != len(tc.err.Error()) {
				if len(tc.err.Error()) > len(err.Error()) {
					assert.Contains(t, tc.err.Error(), err.Error(), "Actual error: %s does not contain expected error: %s", err.Error(), tc.err.Error())
				} else {
					assert.Contains(t, err.Error(), tc.err.Error(), "Actual error: %s does not contain expected error: %s", err.Error(), tc.err.Error())
				}
			} else {
				assert.EqualError(t, err, tc.err.Error())
			}
		} else {
			assert.NoError(t, err)
		}

		server.CloseClientConnections()
		server.Close()
	}
}

func TestSaveBatch(t *testing.T) {
	for name, tc := range SaveBatchTestcases {
		server, mux := newTestServer()
		mux.Handle(tc.path, newMockHandler(tc.method, tc.path, tc.resp, tc.statusCode))
		client, err := NewClient(Config{Host: server.URL}, &http.Client{Transport: &http.Transport{}})
		if err != nil {
			assert.NoError(t, err, "Can not create client: %v", err)
		}
		actualBoInfos, err := client.SaveBatch(tc.bos)
		assert.Equal(t, tc.expectedBoInfos, actualBoInfos, "Unexpected result on testcase '%s': ", name)
		if tc.err != nil && err != nil {
			if len(err.Error()) != len(tc.err.Error()) {
				if len(tc.err.Error()) > len(err.Error()) {
					assert.Contains(t, tc.err.Error(), err.Error(), "Actual error: %s does not contain expected error: %s", err.Error(), tc.err.Error())
				} else {
					assert.Contains(t, err.Error(), tc.err.Error(), "Actual error: %s does not contain expected error: %s", err.Error(), tc.err.Error())
				}
			} else {
				assert.EqualError(t, err, tc.err.Error())
			}
		} else {
			assert.NoError(t, err)
		}

		server.CloseClientConnections()
		server.Close()
	}
}

func TestDeleteByPubID(t *testing.T) {
	for name, tc := range DeleteByPubIDTestcases {
		server, mux := newTestServer()
		mux.Handle(tc.path, newMockHandler(tc.method, tc.path, tc.resp, tc.statusCode))
		client, err := NewClient(Config{Host: server.URL}, &http.Client{Transport: &http.Transport{}})
		if err != nil {
			assert.NoError(t, err, "Can not create client: %v", err)
		}

		errorExpected := tc.err != nil
		gotInfo, err := client.DeleteByPublicID(tc.busobID, tc.publicID)
		if err != nil {
			if !errorExpected {
				t.Fatalf("got unexpected error: %v", err)
				return
			}

			assert.EqualError(t, err, tc.err.Error())
			return
		}

		assert.Equal(t, tc.expectedPubID, gotInfo.PublicID, "Unexpected result on testcase '%s': ", name)
		assert.Equal(t, tc.expectedRecID, gotInfo.RecordID, "Unexpected result on testcase '%s': ", name)
	}
}

func TestDeleteByRecordID(t *testing.T) {
	for name, tc := range DeleteByRecordIDTestcases {
		server, mux := newTestServer()
		mux.Handle(tc.path, newMockHandler(tc.method, tc.path, tc.resp, tc.statusCode))
		client, err := NewClient(Config{Host: server.URL}, &http.Client{Transport: &http.Transport{}})
		if err != nil {
			assert.NoError(t, err, "Can not create client: %v", err)
		}

		actualBoInfo, err := client.DeleteByRecordID(tc.busobID, tc.recordID)
		assert.Equal(t, tc.expectedBoInfo, actualBoInfo, "Unexpected result on testcase '%s': ", name)
		if tc.err != nil {
			if len(err.Error()) != len(tc.err.Error()) {
				if len(err.Error()) > len(tc.err.Error()) {
					assert.Contains(t, err.Error(), tc.err.Error(), "Actual error: %s does not contain expected error: %s", err.Error(), tc.err.Error())
				} else {
					assert.Contains(t, tc.err.Error(), err.Error(), "Actual error: %s does not contain expected error: %s", tc.err.Error(), err.Error())
				}
			} else {
				assert.EqualError(t, err, tc.err.Error())
			}
		} else {
			assert.NoError(t, err)
		}

		server.CloseClientConnections()
		server.Close()
	}
}

func TestDeleteBatch(t *testing.T) {
	for name, tc := range DeleteBatchTestcases {
		server, mux := newTestServer()
		mux.Handle(tc.path, newMockHandler(tc.method, tc.path, tc.resp, tc.statusCode))
		client, err := NewClient(Config{Host: server.URL}, &http.Client{Transport: &http.Transport{}})
		if err != nil {
			assert.NoError(t, err, "Can not create client: %v", err)
		}

		actualBoInfos, err := client.DeleteBatch(tc.bos)
		assert.Equal(t, tc.expectedBoInfos, actualBoInfos, "Unexpected result on testcase '%s': ", name)
		if tc.err != nil && err != nil {
			if len(err.Error()) != len(tc.err.Error()) {
				if len(err.Error()) > len(tc.err.Error()) {
					if len(err.Error()) > len(tc.err.Error()) {
						assert.Contains(t, err.Error(), tc.err.Error(), "Actual error: %s does not contain expected error: %s", err.Error(), tc.err.Error())
					} else {
						assert.Contains(t, tc.err.Error(), err.Error(), "Actual error: %s does not contain expected error: %s", tc.err.Error(), err.Error())
					}
				} else {
					assert.Contains(t, tc.err.Error(), err.Error(), "Actual error: %s does not contain expected error: %s", tc.err.Error(), err.Error())
				}
			} else {
				assert.EqualError(t, err, tc.err.Error())
			}
		} else {
			assert.NoError(t, err)
		}

		server.CloseClientConnections()
		server.Close()
	}
}
