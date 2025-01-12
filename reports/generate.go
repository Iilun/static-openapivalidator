package reports

import (
	"errors"
	"static-openapivalidator/validator"
)

func GenerateReport(results []validator.ValidationResult) (Report, error) {
	var totalRequests, passedRequests, warnRequests, ignoredRequests, totalResponses, warnResponses, passedResponses, ignoredResponses int
	for i := range results {
		switch v := results[i].(type) {
		case *validator.RequestValidationResult:
			totalRequests++
			if v.Status == validator.Success {
				passedRequests++
			} else if v.Status == validator.Warning {
				warnRequests++
			} else if v.Status == validator.Ignored {
				ignoredRequests++
			}
		case *validator.ResponseValidationResult:
			totalResponses++
			if v.Status == validator.Success {
				passedResponses++
			} else if v.Status == validator.Warning {
				warnResponses++
			} else if v.Status == validator.Ignored {
				ignoredResponses++
			}
		default:
			return Report{}, errors.New("got unknown type")
		}
	}
	return Report{
		Summary: Summary{
			TotalRequests:    totalRequests,
			PassedRequests:   passedRequests,
			WarnRequests:     warnRequests,
			IgnoredRequests:  ignoredRequests,
			FailedRequests:   totalRequests - passedRequests - warnRequests - ignoredRequests,
			TotalResponses:   totalResponses,
			PassedResponses:  passedResponses,
			WarnResponses:    warnResponses,
			IgnoredResponses: ignoredResponses,
			FailedResponses:  totalResponses - passedResponses - warnResponses - ignoredResponses,
		},
		Results: results,
	}, nil
}
