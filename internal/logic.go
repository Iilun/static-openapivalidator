package internal

import (
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	struct_validator "github.com/go-playground/validator/v10"
	"github.com/gobwas/glob"
	"gopkg.in/yaml.v3"
	"os"
	"static-openapivalidator/logger"
	test_report "static-openapivalidator/parser"
	"static-openapivalidator/reports"
	"static-openapivalidator/reports/html"
	"static-openapivalidator/reports/json"
	"static-openapivalidator/reports/junit"
	"static-openapivalidator/validator"
)

func (params *Params) Execute() error {
	logger.Enabled = params.Debug

	err := struct_validator.New().Struct(params)
	if err != nil {
		// TODO: improve error messages
		return errors.New("invalid input: " + err.Error())
	}

	logger.Log("Loading configuration from file")
	err = params.loadConfig()
	if err != nil {
		return err
	}

	results, err := params.checkResponses()

	if err != nil {
		return err
	}

	summary, err := params.logResults(results)

	if err != nil {
		return err
	}

	fmt.Println(*summary)
	if summary.FailedResponses > 0 || summary.FailedRequests > 0 {
		return errors.New("run failed")
	}
	return nil
}

func (params *Params) loadConfig() error {
	if params.ConfigFilePath != "" {
		yamlBytes, err := os.ReadFile(params.ConfigFilePath)
		if err != nil {
			return err
		}

		var config Config
		err = yaml.Unmarshal(yamlBytes, &config)
		if err != nil {
			return err
		}
		// Parse globs
		var bannedResponses, bannedRequests, bannedRoutes []glob.Glob

		bannedResponses, err = compileGlobs(config.Ignore.Responses)
		if err != nil {
			return err
		}
		bannedRequests, err = compileGlobs(config.Ignore.Requests)
		if err != nil {
			return err
		}
		bannedRoutes, err = compileGlobs(config.Ignore.Routes)
		if err != nil {
			return err
		}

		params.config = validator.Config{
			IgnoredRequests:  bannedRequests,
			IgnoredResponses: bannedResponses,
			IgnoredRoutes:    bannedRoutes,
			IgnoreServers:    config.Ignore.Servers,
		}
	}
	return nil
}

func compileGlobs(array []string) ([]glob.Glob, error) {
	var result []glob.Glob
	for _, path := range array {
		elem, err := glob.Compile(path)
		if err != nil {
			return nil, err
		}
		result = append(result, elem)
	}
	return result, nil
}

func (params *Params) checkResponses() ([]validator.ValidationResult, error) {
	// Load open api ref
	logger.Log("OpenAPI Spec: loading")
	loader := &openapi3.Loader{Context: params.Ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(params.ApiFilePath)
	if err != nil {
		return nil, err
	}

	logger.Log("OpenAPI Spec: validating")
	// Validate document
	err = doc.Validate(params.Ctx)
	if err != nil {
		return nil, errors.New("openapi validation: " + err.Error())
	}

	if params.config.IgnoreServers {
		// Ignoring servers from spec so requests on any host matches
		logger.Log("OpenAPI Spec: enabling IgnoreServers")
		doc.Servers = openapi3.Servers{}
	}

	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		return nil, err
	}

	// Parse file
	logger.Log("%s: getting parser", params.Format)
	parser, err := getParser(params.Format)
	if err != nil {
		return nil, err
	}

	logger.Log("%s: parsing results from files", params.Format)
	results, err := parser.Parse(params.ReportFilePaths, router, params.config)
	if err != nil {
		return nil, err
	}
	// Validate all request/responses
	openapi3.SchemaErrorDetailsDisabled = true
	return validator.Validate(results, params.Ctx)
}

func (params *Params) logResults(results []validator.ValidationResult) (*reports.Summary, error) {
	var reporters []reports.Reporter

	if params.HtmlFilePath != "" {
		logger.Log("Reporting: enabling HTML reporter")
		reporters = append(reporters, html.NewReporter(params.HtmlFilePath))
	}

	if params.JsonFilePath != "" {
		logger.Log("Reporting: enabling JSON reporter")
		reporters = append(reporters, json.NewReporter(params.JsonFilePath))
	}

	if params.JunitFilePath != "" {
		logger.Log("Reporting: enabling JUNIT reporter")
		reporters = append(reporters, junit.NewReporter(params.JunitFilePath))
	}

	logger.Log("Reporting: generating report on %d results", len(results))
	report, err := reports.GenerateReport(results)
	if err != nil {
		return nil, err
	}
	for i := range reporters {
		err = reporters[i].Generate(report)
		if err != nil {
			return nil, err
		}
	}

	return &report.Summary, nil
}

func getParser(format string) (test_report.Parser, error) {
	switch format {
	case "bruno":
		return test_report.BrunoParser{}, nil
	case "postman":
		return test_report.PostmanParser{}, nil
	default:
		return nil, fmt.Errorf("format %s not supported", format)
	}
}
