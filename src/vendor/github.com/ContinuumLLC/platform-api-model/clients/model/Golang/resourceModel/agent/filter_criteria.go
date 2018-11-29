package agent

import "time"

//Criteria is a type to restrict Filters, and have common understanding between the Resources
type Criteria string

const (
	//Partner : is a Criteria to find all the endpoints for given partner
	Partner Criteria = "Partner"
	//Client : is a Criteria to find all the endpoints for given Partner and Client
	Client Criteria = "Client"
	//Site : is a Criteria to find all the endpoints for given Partner, Client, Site
	Site Criteria = "Site"
	//Endpoint : is a Criteria to find all the endpoints for given Partner, Client, Site and Endpoint
	Endpoint Criteria = "Endpoint"
	//URL : is a Criteria to find all the endpoints for given URL
	URL Criteria = "URL"
)

//FilterCriteria is a struct to hold filter criteria for Manifest and Mailbox message
type FilterCriteria struct {
	Filter  Criteria                `json:"filter,omitempty"`
	Mapping []FilterEndpointMapping `json:"mapping,omitempty"`
	URL     string                  `json:"url,omitempty"`
}

//FilterEndpointMapping is a struct to hold filter endpoint Mapping
type FilterEndpointMapping struct {
	Partner        string    `json:"partner,omitempty"`
	Client         string    `json:"client,omitempty"`
	Site           string    `json:"site,omitempty"`
	Endpoints      []string  `json:"endpoints,omitempty"`
	DCCreatedTSUTC time.Time `json:"dcCreatedTSUTC,omitempty"`
}

//FilterCriteriaSchema is schema to validate payload
var FilterCriteriaSchema = `{
	"$schema": "http://json-schema.org/draft-04/schema#",
	"title": "EnableManifestPost",
	"description": "This schema represents the JSON body of requests for POST enableManifest",
	"type": "object",
	"properties": {
		"filter": {
			"type": "string"
		},
		"mapping": {
			"type": "array",
			"items": {
				"$ref": "#/definitions/mappingDef"
			},
			"minItems": 1
		}
	},
	"required": [
		"filter",
		"mapping"
	],
	"minItems": 1,
	"definitions": {
		"mappingDef": {
			"type": "object",
			"properties": {
				"partner": {
					"type": "string"
				},
				"client": {
					"type": "string"
				},
				"site": {
					"type": "string"
				},
				"endpoints": {
					"type": "array",
					"items": {
						"type": "string"
					},
					"uniqueItems": true
				}
			},
			"required": [
				"partner"
			]
		}
	}
}`
