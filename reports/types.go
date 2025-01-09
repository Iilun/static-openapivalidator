package reports

import "static-openapivalidator/validator"

type Report struct {
	Summary Summary                      `json:"summary"`
	Results []validator.ValidationResult `json:"results"`
}

type Summary struct {
	TotalRequests   int `json:"totalRequests"`
	PassedRequests  int `json:"passedRequests"`
	FailedRequests  int `json:"failedRequests"`
	TotalResponses  int `json:"totalResponses"`
	PassedResponses int `json:"passedResponses"`
	FailedResponses int `json:"failedResponses"`
}
