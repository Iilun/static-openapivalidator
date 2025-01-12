package bruno

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

type Parser struct{}

func (p Parser) Parse(reportFilePaths []string, router routers.Router, config validator.Config) ([]validator.TestResult, error) {
	var results []Result

	for _, path := range reportFilePaths {
		var reports []Report
		reportBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(reportBytes, &reports)
		if err != nil {
			return nil, errors.New(path + ": " + err.Error())
		}

		if len(reports) == 0 {
			return nil, errors.New(path + ": no report in report file")
		}

		for i := range reports {
			for j := range reports[i].Results {
				if len(reportFilePaths) > 1 {
					reports[i].Results[j].FileOrigin = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
				}
				results = append(results, reports[i].Results[j])
			}
		}
	}

	translated, err := translateResults(results, router, config)
	if err != nil {
		return nil, err
	}
	return translated, nil
}

func translateResults(results []Result, router routers.Router, config validator.Config) ([]validator.TestResult, error) {
	var final []validator.TestResult
	for i := range results {
		res, err := brunoToOpenAPI(results[i], router, config)
		if err != nil {
			return nil, err
		}
		final = addResultToArray(final, res, config)
	}
	return final, nil
}

func addResultToArray(array []validator.TestResult, res validator.TestResult, config validator.Config) []validator.TestResult {
	// Check if request is ignored
	for _, path := range config.BannedRequests {
		if path.Match(res.Id) {
			res.Request.Ignored = true
		}
	}

	// Check if request is ignored
	for _, path := range config.BannedResponses {
		if path.Match(res.Id) {
			res.Response.Ignored = true
		}
	}

	// Route was ignored or both were ignored
	if res.Request.Ignored && res.Response.Ignored {
		return array
	}

	return append(array, res)
}

func brunoToOpenAPI(result Result, router routers.Router, config validator.Config) (validator.TestResult, error) {
	request, err := translateRequest(result.Request, router, config)
	if err != nil {
		return validator.TestResult{}, err
	}
	response, err := translateResponse(result.Response, request)
	if err != nil {
		return validator.TestResult{}, err
	}
	return validator.TestResult{
		Request:  request,
		Response: response,
		Id:       formatId(result.FileOrigin, result.Test.Filename),
	}, nil
}

func formatId(fileOrigin, filename string) string {
	filename = strings.TrimSuffix(filename, ".bru")
	filename = strings.TrimSuffix(filename, "-muted-")
	filename = strings.TrimSpace(filename)
	if fileOrigin != "" {
		filename = fileOrigin + "/" + filename
	}
	return filename
}

func getHeaderValue(header string, headers map[string]any) string {
	for key, value := range headers {
		if strings.EqualFold(key, header) {
			return fmt.Sprintf("%s", value)
		}
	}
	return ""
}

func translateRequest(brunoRequest Request, router routers.Router, config validator.Config) (*validator.TestRequest, error) {
	// Translate request
	var requestBody io.Reader
	var prettyJSON bytes.Buffer
	var parsingError string
	if brunoRequest.Body != "" {
		// Multipart body
		if strings.Contains(getHeaderValue("Content-Type", brunoRequest.Headers), "multipart/form-data") {
			// Multipart in JSON report has no data as to what was sent, and not event the set fields
			parsingError = "multipart/form-data is not supported"
		} else {
			// JSON Body
			requestBody = strings.NewReader(string(brunoRequest.Body))
			if err := json.Indent(&prettyJSON, []byte(brunoRequest.Body), "", "  "); err != nil {
				return nil, err
			}
		}
	}

	// Bruno does not escape URLs
	splitted := strings.Split(strings.Split(brunoRequest.Url, "?")[0], "/")
	for i := range splitted {
		splitted[i] = url.PathEscape(splitted[i])
	}

	parsedUrl, err := url.Parse(strings.Join(splitted, "/"))
	if err != nil {
		return nil, err
	}

	ignored := false
	for _, path := range config.BannedRoutes {
		if path.Match(parsedUrl.Path) {
			ignored = true
		}
	}

	httpReq, err := http.NewRequest(brunoRequest.Method, parsedUrl.String(), requestBody)
	if err != nil {
		return nil, err
	}
	// Set headers
	for header, value := range brunoRequest.Headers {
		httpReq.Header.Set(header, fmt.Sprintf("%s", value))
	}

	route, pathParams, err := router.FindRoute(httpReq)
	if err != nil {
		if errors.Is(err, routers.ErrPathNotFound) {
			parsingError = fmt.Sprintf("could not find route for %s %s: %v", brunoRequest.Method, parsedUrl.String(), err)
		} else if errors.Is(err, routers.ErrMethodNotAllowed) {
			parsingError = fmt.Sprintf("bad method for %s %s: %v", brunoRequest.Method, parsedUrl.String(), err)
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

func translateResponse(brunoResponse Response, request *validator.TestRequest) (*validator.TestResponse, error) {
	headers := http.Header{}
	for header, value := range brunoResponse.Headers {
		headers.Set(header, fmt.Sprintf("%s", value))
	}
	var bodyReader io.ReadCloser
	var prettyJSON bytes.Buffer
	if brunoResponse.Body != "" {
		bodyReader = io.NopCloser(strings.NewReader(string(brunoResponse.Body)))
		if err := json.Indent(&prettyJSON, []byte(brunoResponse.Body), "", "  "); err != nil {
			return nil, err
		}
	}
	var parsingError string
	if request.Route == nil {
		parsingError = "no route found"
	}
	return &validator.TestResponse{
		ResponseValidationInput: &openapi3filter.ResponseValidationInput{
			RequestValidationInput: request.RequestValidationInput,
			Status:                 brunoResponse.Status,
			Header:                 headers,
			Body:                   bodyReader,
		},
		Body:         prettyJSON.String(),
		ParsingError: parsingError,
		Ignored:      request.Ignored,
	}, nil
}
