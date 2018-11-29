package cherwell

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {

	type findTestcases struct {
		resp          string
		expected      []BusinessObject
		mockedID      string
		expectedError error
		statusCode    int
		operator string
	}

	testcases := map[string]findTestcases{
		"Find with all fields": {
			resp: `{
		"businessObjects": [
		 {
			"busObId": "6d",
			"busObPublicId": "19",
			"busObRecId": "94",
			"fields": [
			  {
				"dirty": false,
				"displayName": "Service Order Number",
				"fieldId": "1",
				"html": null,
				"name": "CartItemID",
				"value": ""
			  },
			  {
				"dirty": false,
				"displayName": "Description",
				"fieldId": "2",
				"html": null,
				"name": "Incident description",
				"value": ""
			  }
			],
			"links": [],
			"errorCode": null,
			"errorMessage": null,
			"hasError": false
		 },
		 {
			"busObId": "7d",
			"busObPublicId": "23",
			"busObRecId": "78",
			"fields": [
			  {
				"dirty": false,
				"displayName": "Service Order Number",
				"fieldId": "1",
				"html": null,
				"name": "CartItemID",
				"value": ""
			  },
			  {
				"dirty": false,
				"displayName": "Description",
				"fieldId": "2",
				"html": null,
				"name": "Incident description",
				"value": ""
  	         }
			],
			"links": [],
			"errorCode": null,
			"errorMessage": null,
			"hasError": false
		 }
		],
		"hasPrompts": false,
		"links": [],
		"prompts": [],
		"searchResultsFields": [],
		"simpleResults": null,
		"totalRows": 1,
		"errorCode": null,
		"errorMessage": null,
		"hasError": false
	  }`,
			expected: []BusinessObject{
				{
					BusinessObjectInfo: BusinessObjectInfo{
						ID:       "6d",
						RecordID: "94",
						PublicID: "19",
					},
					Fields: []FieldTemplateItem{
						{
							DisplayName: "Service Order Number",
							FieldID:     "1",
							Name:        "CartItemID",
						},
						{
							DisplayName: "Description",
							FieldID:     "2",
							Name:        "Incident description",
						},
					},
				},
				{
					BusinessObjectInfo: BusinessObjectInfo{
						ID:       "7d",
						RecordID: "78",
						PublicID: "23",
					},
					Fields: []FieldTemplateItem{
						{
							DisplayName: "Service Order Number",
							FieldID:     "1",
							Name:        "CartItemID",
						},
						{
							DisplayName: "Description",
							FieldID:     "2",
							Name:        "Incident description",
						},
					},
				},
			},
			mockedID:   "123",
			statusCode: 200,
			operator: "eq",
		},

		"Find with err in response": {
			resp: `{
		"businessObjects": null,
		"hasPrompts": false,
		"links": [],
		"prompts": [],
		"searchResultsFields": [],
		"simpleResults": null,
		"totalRows": 1,
		"errorCode": "",
		"errorMessage": "",
		"hasError": true
	  }`,
			expected:      nil,
			mockedID:      "",
			expectedError: errors.New("BAD_REQUEST"),
			statusCode:    400,
			operator: "eq",
		},

		"Find with error in performRequest method": {
			expected:      nil,
			mockedID:      "123",
			expectedError: errors.New("INTERNAL_SERVER_ERROR"),
			statusCode:    500,
			operator: "eq",
		},

		"Find with error in resp.Responses": {
			resp: `{
		"businessObjects": [
		{
			"busObId": "6d",
			"busObPublicId": "19",
			"busObRecId": "94",
			"fields": [
			  {
				"dirty": false,
				"displayName": "Service Order Number",
				"fieldId": "1",
				"html": null,
				"name": "CartItemID",
				"value": ""
			  }
			],
			"links": [],
			"errorCode": "",
			"errorMessage": "some error",
			"hasError": true
		}
		],
		"hasPrompts": false,
		"links": [],
		"prompts": [],
		"searchResultsFields": [],
		"simpleResults": null,
		"totalRows": 1,
		"errorCode": null,
		"errorMessage": null,
		"hasError": false
	  }`,
			mockedID:      "123",
			statusCode:    400,
			expectedError: errors.New("BAD_REQUEST"),
			operator: "eq",
		},

		"Find with invalid operator": {
			resp: `{
		"businessObjects": null,
		"hasPrompts": false,
		"links": [],
		"prompts": [],
		"searchResultsFields": [],
		"simpleResults": null,
		"totalRows": 1,
		"errorCode": "",
		"errorMessage": "",
		"hasError": true
	  }`,
			expected:      nil,
			mockedID:      "",
			expectedError: errors.New("BAD_REQUEST"),
			statusCode:    400,
			operator: "invalid",
		},
	}

	for _, tc := range testcases {
		server, mux := newTestServer()
		mockHandler := newMockHandler(http.MethodPost, searchEndpoint, tc.resp, tc.statusCode)

		req := NewSearchRequest(tc.mockedID)
		req.AddFilter("1", tc.operator, "fieldvalue1")
		req.AddFilter("2", tc.operator, "fieldvalue2")

		mux.Handle(searchEndpoint, mockHandler)
		client, err := NewClient(Config{Host: server.URL}, &http.Client{Transport: &http.Transport{}})
		assert.NoError(t, err, "Can not create client: %v", err)

		bos, err := client.Find(*req)
		if tc.expectedError == nil {
			assert.NoError(t, err)
		}
		assert.Equal(t, tc.expected, bos)
	}
}

func TestFindBoInfos(t *testing.T) {
	server, mux := newTestServer()

	type findTestcases struct {
		resp     string
		expected []BusinessObject
	}

	testcases := map[string]findTestcases{
		"Find with all fields": {
			resp: `{
		"businessObjects": [
		  {
			"busObId": "6d",
			"busObPublicId": "19",
			"busObRecId": "94",
			"links": [],
			"errorCode": null,
			"errorMessage": null,
			"hasError": false
		  },
		  {
			"busObId": "7d",
			"busObPublicId": "23",
			"busObRecId": "78",
			"links": [],
			"errorCode": null,
			"errorMessage": null,
			"hasError": false
		  }
		],
		"hasPrompts": false,
		"links": [],
		"prompts": [],
		"searchResultsFields": [],
		"simpleResults": null,
		"totalRows": 1,
		"errorCode": null,
		"errorMessage": null,
		"hasError": false
	  }`,
			expected: []BusinessObject{
				{
					BusinessObjectInfo: BusinessObjectInfo{
						ID:       "6d",
						RecordID: "94",
						PublicID: "19",
					},
				},
				{
					BusinessObjectInfo: BusinessObjectInfo{
						ID:       "7d",
						RecordID: "78",
						PublicID: "23",
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		mockHandler := newMockHandler(http.MethodPost, searchEndpoint, tc.resp, http.StatusOK)

		req := NewSearchRequest("123")
		req.AddFilter("1", " eq", "fieldvalue1")
		req.AddFilter("2", "eq", "fieldvalue2")

		mux.Handle(searchEndpoint, mockHandler)
		client, err := NewClient(Config{Host: server.URL}, &http.Client{Transport: &http.Transport{}})
		assert.NoError(t, err, "Can not create client: %v", err)

		bos, err := client.FindBoInfos(*req)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, bos)
	}
}

func TestSetSpecificFields(t *testing.T) {

	type findTestcase struct {
		input       SearchRequest
		expected    SearchRequest
		inputFields []string
	}

	tc := findTestcase{
		input: SearchRequest{
			IncludeAllFields: true,
		},
		expected: SearchRequest{
			Fields: []string{
				"1", "2", "3",
			},
			IncludeAllFields: false,
		},
		inputFields: []string{
			"1", "2", "3",
		},
	}

	tc.input.SetSpecificFields(tc.inputFields)
	assert.Equal(t, tc.expected, tc.input)
}

func TestAppendSpecificFields(t *testing.T) {

	type findTestcase struct {
		input        SearchRequest
		expected     SearchRequest
		appendFields []string
	}

	tc := findTestcase{
		input: SearchRequest{
			Fields: []string{
				"1",
			},
			IncludeAllFields: true,
		},
		expected: SearchRequest{
			Fields: []string{
				"1", "2", "3",
			},
			IncludeAllFields: false,
		},
		appendFields: []string{
			"2", "3",
		},
	}

	tc.input.AppendSpecificFields(tc.appendFields...)
	assert.Equal(t, tc.expected, tc.input)
}
