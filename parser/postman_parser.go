package test_report

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"static-openapivalidator/validator"
	"strings"
)

type PostmanParser struct{}

func (p PostmanParser) Parse(reportFilePaths []string, router routers.Router, config validator.Config) ([]validator.TestResult, error) {
	var results []PostmanExecution

	for _, path := range reportFilePaths {
		var report PostmanReport
		reportBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(reportBytes, &report)
		if err != nil {
			return nil, errors.New(path + ": " + err.Error())
		}

		for i := range report.Run.Executions {
			if len(reportFilePaths) > 1 {
				report.Run.Executions[i].FileOrigin = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
			}
			testId := report.Run.Executions[i].Id
			// Find the testId in the report path

			report.Run.Executions[i].JsonPath = findPathToId(testId, "", report.Collection.Item)
			results = append(results, report.Run.Executions[i])
		}
	}
	translated, err := translatePostmanResults(results, router, config)
	if err != nil {
		return nil, err
	}

	return translated, nil
}

func findPathToId(id, path string, items []PostmanItem) string {
	for i := range items {
		if id == items[i].Id {
			return path
		}
		pathInChildren := findPathToId(id, items[i].Name, items[i].Item)
		if pathInChildren != "" {
			if path != "" {
				return path + "/" + pathInChildren
			}
			return pathInChildren
		}
	}
	return ""
}

func translatePostmanResults(results []PostmanExecution, router routers.Router, config validator.Config) ([]validator.TestResult, error) {
	var final []validator.TestResult
	for i := range results {
		res, err := postmanToOpenAPI(results[i], router, config)
		if err != nil {
			return nil, err
		}
		final = addResultToArray(final, res, config)
	}
	return final, nil
}

func postmanToOpenAPI(result PostmanExecution, router routers.Router, config validator.Config) (validator.TestResult, error) {
	request, err := translatePostmanRequest(result.Request, router, config)
	if err != nil {
		return validator.TestResult{}, err
	}
	response, err := translatePostmanResponse(result.Response, request)
	if err != nil {
		return validator.TestResult{}, err
	}
	return validator.TestResult{
		Request:  request,
		Response: response,
		Id:       formatPostmanId(result),
	}, nil
}

func formatPostmanId(result PostmanExecution) string {
	var final []string
	for _, elem := range []string{result.FileOrigin, result.JsonPath, result.Item.Name} {
		if elem != "" {
			final = append(final, elem)
		}
	}
	return strings.Join(final, "/")
}

func translatePostmanRequest(postmanRequest PostmanRequest, router routers.Router, config validator.Config) (*validator.TestRequest, error) {
	// Translate request
	var requestBody io.Reader
	var prettyJSON bytes.Buffer
	var parsingError string

	if postmanRequest.Body.Raw != "" {
		// Multipart body
		//if strings.Contains(getHeaderValue("Content-Type", postmanRequest.Header), "multipart/form-data") {
		//	// Multipart in JSON report has no data as to what was sent, and not event the set fields
		//	parsingError = "multipart/form-data is not supported"
		//} else {
		// JSON Body
		requestBody = strings.NewReader(postmanRequest.Body.Raw)
		if err := json.Indent(&prettyJSON, []byte(postmanRequest.Body.Raw), "", "  "); err != nil {
			return nil, errors.New("could not format request body: " + err.Error())
		}
	}

	parsedUrl, err := url.Parse(postmanRequest.URL.GetUrl())
	if err != nil {
		return nil, err
	}

	ignored := false
	for _, path := range config.IgnoredRoutes {
		if path.Match(parsedUrl.Path) {
			ignored = true
		}
	}

	httpReq, err := http.NewRequest(postmanRequest.Method, parsedUrl.String(), requestBody)
	if err != nil {
		return nil, err
	}

	// Set headers
	for _, header := range postmanRequest.Header {
		httpReq.Header.Set(header.Key, header.Value)
	}

	// TODO: move this out as it should be common
	route, pathParams, err := router.FindRoute(httpReq)
	if err != nil {
		if errors.Is(err, routers.ErrPathNotFound) {
			parsingError = fmt.Sprintf("could not find route for %s %s: %v", postmanRequest.Method, parsedUrl.String(), err)
		} else if errors.Is(err, routers.ErrMethodNotAllowed) {
			parsingError = fmt.Sprintf("bad method for %s %s: %v", postmanRequest.Method, parsedUrl.String(), err)
		} else {
			return nil, err
		}
	} else {
		// Disabling security checks
		route.Spec.Security = nil
	}

	request := validator.TestRequest{
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request:    httpReq,
			PathParams: pathParams,
			Route:      route,
		},
		Body:         prettyJSON.String(),
		ParsingError: parsingError,
		Ignored:      ignored,
	}

	return &request, nil
}

func getPostmanHeaderValue(headers []PostmanHeader, headerName string) string {
	for _, header := range headers {
		if strings.EqualFold(headerName, header.Key) {
			return header.Value
		}
	}
	return ""
}

func translatePostmanResponse(postmanResponse PostmanResponse, request *validator.TestRequest) (*validator.TestResponse, error) {
	headers := http.Header{}
	for _, header := range postmanResponse.Header {
		headers.Set(header.Key, header.Value)
	}
	var bodyReader io.ReadCloser
	var prettyJSON bytes.Buffer
	if postmanResponse.Stream != "" {
		bodyReader = io.NopCloser(strings.NewReader(string(postmanResponse.Stream)))

		contentType := getPostmanHeaderValue(postmanResponse.Header, "Content-Type")

		if strings.Contains(contentType, "json") {
			if err := json.Indent(&prettyJSON, []byte(postmanResponse.Stream), "", "  "); err != nil {
				return nil, errors.New("could not format response body: " + err.Error())
			}
		} else {
			prettyJSON.WriteString(string(postmanResponse.Stream))
		}
	}
	var parsingError string
	if request.Route == nil {
		parsingError = "no route found"
	}
	return &validator.TestResponse{
		ResponseValidationInput: &openapi3filter.ResponseValidationInput{
			RequestValidationInput: request.RequestValidationInput,
			Status:                 postmanResponse.Code,
			Header:                 headers,
			Body:                   bodyReader,
		},
		Body:         prettyJSON.String(),
		ParsingError: parsingError,
		Ignored:      request.Ignored,
	}, nil
}
