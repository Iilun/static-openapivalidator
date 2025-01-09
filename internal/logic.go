package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"os"
	test_report "static-openapivalidator/parser"
	"static-openapivalidator/parser/bruno"
	"static-openapivalidator/reports"
	"static-openapivalidator/reports/html"
	"static-openapivalidator/reports/json"
	"static-openapivalidator/reports/junit"
	"static-openapivalidator/validator"
)

type Params struct {
	Ctx            context.Context
	ApiFilePath    string
	ReportFilePath string
	Format         string
	JunitFilePath  string
	HtmlFilePath   string
	JsonFilePath   string
}

func (params Params) Execute() error {
	// TODO: strcut validation
	results, err := params.checkResponses()

	if err != nil {
		return err
	}

	return params.logResults(results)
}

func (params Params) checkResponses() ([]validator.ValidationResult, error) {
	// Read report file
	reportBytes, err := os.ReadFile(params.ReportFilePath)
	if err != nil {
		return nil, err
	}

	// Load open api ref
	loader := &openapi3.Loader{Context: params.Ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(params.ApiFilePath)
	if err != nil {
		return nil, err
	}

	// Validate document
	err = doc.Validate(params.Ctx)
	if err != nil {
		return nil, errors.New("openapi validation: " + err.Error())
	}

	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		return nil, err
	}

	// Parse file
	parser, err := getParser(params.Format)
	if err != nil {
		return nil, err
	}

	results, err := parser.Parse(reportBytes, router)
	if err != nil {
		return nil, err
	}
	// Validate all request/responses
	openapi3.SchemaErrorDetailsDisabled = true
	return validator.Validate(results, params.Ctx)
}

func (params Params) logResults(results []validator.ValidationResult) error {
	var reporters []reports.Reporter

	if params.HtmlFilePath != "" {
		reporters = append(reporters, html.NewReporter(params.HtmlFilePath))
	}

	if params.JsonFilePath != "" {
		reporters = append(reporters, json.NewReporter(params.JsonFilePath))
	}

	if params.JunitFilePath != "" {
		reporters = append(reporters, junit.NewReporter(params.JunitFilePath))
	}

	if len(reporters) > 0 {
		report, err := reports.GenerateReport(results)
		if err != nil {
			return err
		}
		for i := range reporters {
			err = reporters[i].Generate(report)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getParser(format string) (test_report.Parser, error) {
	switch format {
	case "bruno":
		return bruno.Parser{}, nil
	default:
		return nil, fmt.Errorf("Format %s not supported", format)
	}
}
