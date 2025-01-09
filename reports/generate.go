package reports

import (
	"errors"
	"static-openapivalidator/validator"
)

func GenerateReport(results []validator.ValidationResult) (Report, error) {
	var totalRequests, passedRequests, warnRequests, totalResponses, warnResponses, passedResponses int
	for i := range results {
		switch v := results[i].(type) {
		case *validator.RequestValidationResult:
			totalRequests++
			if v.Status == validator.Success {
				passedRequests++
			} else if v.Status == validator.Warning {
				warnRequests++
			}
		case *validator.ResponseValidationResult:
			totalResponses++
			if v.Status == validator.Success {
				passedResponses++
			} else if v.Status == validator.Warning {
				warnResponses++
			}
		default:
			return Report{}, errors.New("got unknown type")
		}
	}

	return Report{
		Summary: Summary{
			TotalRequests:   totalRequests,
			PassedRequests:  passedRequests,
			WarnRequests:    warnRequests,
			FailedRequests:  totalRequests - passedRequests - warnRequests,
			TotalResponses:  totalResponses,
			PassedResponses: passedResponses,
			WarnResponses:   warnResponses,
			FailedResponses: totalResponses - passedResponses - warnResponses,
		},
		Results: results,
	}, nil
}
