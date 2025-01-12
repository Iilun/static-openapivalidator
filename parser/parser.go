package test_report

import (
	"github.com/getkin/kin-openapi/routers"
	"static-openapivalidator/validator"
)

// Implement parser to parse a file
type Parser interface {
	Parse(reportFilePaths []string, router routers.Router, config validator.Config) ([]validator.TestResult, error)
}
