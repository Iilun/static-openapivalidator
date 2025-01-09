package reports

import "static-openapivalidator/validator"

type Report struct {
	Summary Summary                      `json:"summary"`
	Results []validator.ValidationResult `json:"results"`
}

type Summary struct {
	TotalRequests   int `json:"totalRequests"`
	PassedRequests  int `json:"passedRequests"`
	WarnRequests    int `json:"warnRequests"`
	FailedRequests  int `json:"failedRequests"`
	TotalResponses  int `json:"totalResponses"`
	PassedResponses int `json:"passedResponses"`
	WarnResponses   int `json:"warnResponses"`
	FailedResponses int `json:"failedResponses"`
}
