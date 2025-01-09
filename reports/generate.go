package reports

import (
	"errors"
	"static-openapivalidator/validator"
)

func GenerateReport(results []validator.ValidationResult) (Report, error) {
	var totalRequests, passedRequests, totalResponses, passedResponses int
	for i := range results {
		switch v := results[i].(type) {
		case *validator.RequestValidationResult:
			totalRequests++
			if v.Status == "success" {
				passedRequests++
			}
		case *validator.ResponseValidationResult:
			totalResponses++
			if v.Status == "success" {
				passedResponses++
			}
		default:
			return Report{}, errors.New("got unknown type")
		}
	}

	return Report{
		Summary: Summary{
			TotalRequests:   totalRequests,
			PassedRequests:  passedRequests,
			FailedRequests:  totalRequests - passedRequests,
			TotalResponses:  totalResponses,
			PassedResponses: passedResponses,
			FailedResponses: totalResponses - passedResponses,
		},
		Results: results,
	}, nil
}
