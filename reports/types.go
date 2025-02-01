package reports

import (
	"fmt"
	"static-openapivalidator/validator"
)

type Report struct {
	Summary Summary                      `json:"summary"`
	Results []validator.ValidationResult `json:"results"`
}

type Summary struct {
	TotalRequests    int `json:"totalRequests"`
	PassedRequests   int `json:"passedRequests"`
	WarnRequests     int `json:"warnRequests"`
	FailedRequests   int `json:"failedRequests"`
	IgnoredRequests  int `json:"ignoredRequests"`
	TotalResponses   int `json:"totalResponses"`
	PassedResponses  int `json:"passedResponses"`
	WarnResponses    int `json:"warnResponses"`
	FailedResponses  int `json:"failedResponses"`
	IgnoredResponses int `json:"ignoredResponses"`
}

func (s Summary) String() string {
	return fmt.Sprintf(`Total requests: %d
Passed requests: %d
Warn requests: %d
Failed requests: %d
Ignored requests: %d
Total responses: %d
Passed respones: %d
Warn responses: %d
Failed reponses: %d
Ignored responses: %d`,
		s.TotalRequests,
		s.PassedRequests,
		s.WarnRequests,
		s.FailedRequests,
		s.IgnoredRequests,
		s.TotalResponses,
		s.PassedResponses,
		s.WarnResponses,
		s.FailedResponses,
		s.IgnoredResponses)
}
