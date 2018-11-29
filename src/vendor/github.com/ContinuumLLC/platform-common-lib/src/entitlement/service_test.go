package entitlement

import (
	"fmt"
	apiModel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/entitlement"
	"github.com/gocql/gocql"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

var (
	emptyFeatures    = []apiModel.Feature{}
	cacheSize        = 104857600
	entitlementMsURL = "http://localhost:8888/entitlement/v1"
	cacheDataTTLSec  = 600
	serviceInstance  Service
)

func init() {
	serviceInstance = NewEntitlementService(http.DefaultClient, entitlementMsURL, cacheDataTTLSec, cacheSize)
}

func TestGetPartnerFeatures(t *testing.T) {
	partner := "1d4400c0"
	partner2 := "222222"
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerResponder(partner, emptyFeatures, t)
	features, err := serviceInstance.GetPartnerFeatures(partner)
	if err != nil {
		t.Fatalf("Got an error: %v", err)
	}
	if len(features) != 0 {
		t.Fatalf("Expected empty features for partner %v, but got %v", partner, features)
	}

	httpmock.DeactivateAndReset()
	httpmock.Activate()
	payload := []apiModel.Feature{
		{
			Name: "TASKING_BASIC",
			ID:   gocql.TimeUUID().String(),
		},
	}
	registerResponder(partner2, payload, t)

	filledFeatures, err := serviceInstance.GetPartnerFeatures(partner2)
	if err != nil {
		t.Fatalf("Got an error: %v", err)
	}
	if len(filledFeatures) != 1 {
		t.Fatalf("Expected 1 feature for partner %v, but got %v", partner, len(filledFeatures))
	}

	// Try to get features from cache
	httpmock.DeactivateAndReset()
	cachedFeatures, err := serviceInstance.GetPartnerFeatures(partner2)
	if err != nil {
		t.Fatalf("Got an error: %v", err)
	}
	if len(cachedFeatures) != 1 {
		t.Fatalf("Expected 1 feature for partner %v, but got %v", partner, len(cachedFeatures))
	}
}

func TestIsPartnerAuthorized(t *testing.T) {
	var (
		testFeature = "TASKING_TEST"
		testCases   = []struct {
			name              string
			partner           string
			payload           []apiModel.Feature
			expectedResult    bool
			registerResponder bool
		}{
			{
				name:              "No response from Entitlement",
				partner:           "partnerNotAuthorized",
				payload:           nil,
				expectedResult:    false,
				registerResponder: false,
			},
			{
				name:              "Partner Not Authorized",
				partner:           "partnerNotAuthorized",
				payload:           nil,
				expectedResult:    false,
				registerResponder: true,
			},
			{
				name:    "Partner Authorized",
				partner: "partnerAuthorized",
				payload: []apiModel.Feature{
					{
						Name: testFeature,
						ID:   gocql.TimeUUID().String(),
					},
				},
				expectedResult:    true,
				registerResponder: true,
			},
		}
	)
	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			if tc.registerResponder {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()
				registerResponder(tc.partner, tc.payload, t)
			}
			isPartnerInACL := serviceInstance.IsPartnerAuthorized(tc.partner, testFeature)
			fmt.Printf("%s: %v\n", tc.name, tc.payload)
			if isPartnerInACL != tc.expectedResult {
				t.Fatalf("\nTestCase[%s]: Partner is in ACL = %t, but want %t", tc.name, isPartnerInACL, tc.expectedResult)
			}
		})
	}
}

func registerResponder(partnerID string, payload []apiModel.Feature, t *testing.T) {

	entitlementURL := fmt.Sprintf("%s/partners/%s/features", entitlementMsURL, partnerID)

	t.Logf("Registered HTTP responder on URL: %s, payload: %v\n", entitlementURL, payload)
	httpmock.RegisterResponder("GET", entitlementURL,
		func(req *http.Request) (*http.Response, error) {
			if len(payload) == 0 {
				return httpmock.NewJsonResponse(http.StatusNoContent, payload)
			}

			return httpmock.NewJsonResponse(http.StatusOK, payload)
		},
	)
}
