package test_report

import "static-openapivalidator/validator"

func addResultToArray(array []validator.TestResult, res validator.TestResult, config validator.Config) []validator.TestResult {
	// Check if request is ignored
	for _, path := range config.IgnoredRequests {
		if path.Match(res.Id) {
			res.Request.Ignored = true
		}
	}

	// Check if request is ignored
	for _, path := range config.IgnoredResponses {
		if path.Match(res.Id) {
			res.Response.Ignored = true
		}
	}

	return append(array, res)
}
