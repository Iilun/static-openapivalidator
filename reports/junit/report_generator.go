package junit

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	junit_xml "github.com/jstemmer/go-junit-report/v2/junit"
	"os"
	"static-openapivalidator/reports"
	"static-openapivalidator/validator"
	"strings"
)

func NewReporter(filePath string) Reporter {
	return Reporter{filePath: filePath}
}

type Reporter struct {
	filePath string
}

func (r Reporter) Generate(report reports.Report) error {
	// Group results by URL

	groups := make(map[string][]validator.ValidationResult)

	for i := range report.Results {
		url := report.Results[i].GetUrl()
		if len(groups[url]) == 0 {
			groups[url] = []validator.ValidationResult{report.Results[i]}
		} else {
			groups[url] = append(groups[url], report.Results[i])
		}
	}

	var suites junit_xml.Testsuites
	for url, tests := range groups {
		suite := junit_xml.Testsuite{
			Name: url,
			ID:   len(tests),
		}
		for _, test := range tests {
			suite.AddTestcase(createTestcaseForTest(url, test))
		}
		suites.AddSuite(suite)
	}

	bytes, err := xml.Marshal(suites)
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, bytes, 0644)
}

func createTestcaseForTest(url string, test validator.ValidationResult) junit_xml.Testcase {

	testName := fmt.Sprintf("%s - %s", test.GetTestId(), test.GetType())

	tc := junit_xml.Testcase{
		Classname: url,
		Name:      testName,
	}

	if test.GetErrorSummary() != "" {
		tc.Failure = &junit_xml.Result{
			Message: "Failed",
			Data:    formatOutput(test),
		}
	} else {
		tc.SystemOut = &junit_xml.Output{Data: formatOutput(test)}
	}
	return tc
}

func formatOutput(test validator.ValidationResult) string {

	var sb strings.Builder
	// Always log body and headers
	if len(test.GetHeaders()) > 0 {
		sb.WriteString("Headers:\n")
		for key, values := range test.GetHeaders() {
			sb.WriteString(fmt.Sprintf("\t%s: %v", key, values))
		}
		sb.WriteString("\n")
	}

	if test.GetBody() != "" {
		sb.WriteString("Body:\n")
		sb.WriteString(indent(test.GetBody(), "\t"))
		sb.WriteString("\n")
	}

	// Log errors if exists

	if test.GetErrorSummary() != "" {
		sb.WriteString(fmt.Sprintf("Error summary: %s \n", test.GetErrorSummary()))
	}

	if len(test.GetErrors()) > 0 {
		for i := range test.GetErrors() {
			sb.WriteString(fmt.Sprintf("Error #%d:\n", i+1))
			sb.WriteString(test.GetErrors()[i].Title)
			if test.GetErrors()[i].Schema != "" {
				sb.WriteString("\n\tSchema:\n")
				sb.WriteString(indent(test.GetErrors()[i].Schema, "\t"))
			}
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func indent(s, indent string) string {
	return indent + strings.ReplaceAll(s, "\n", "\n"+indent)
}
