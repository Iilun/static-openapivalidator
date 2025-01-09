package bruno

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"io"
	"net/http"
	"static-openapivalidator/validator"
	"strings"
)

type Parser struct{}

func (p Parser) Parse(report []byte, router routers.Router) ([]validator.TestResult, error) {
	var reports []Report
	err := json.Unmarshal(report, &reports)
	if err != nil {
		return nil, err
	}

	if len(reports) == 0 {
		return nil, errors.New("no report in report file")
	}

	var results []Result
	for i := range reports {
		results = append(results, reports[i].Results...)
	}

	translated, err := translateResults(results, router)
	if err != nil {
		return nil, err
	}
	return translated, nil
}

func translateResults(results []Result, router routers.Router) ([]validator.TestResult, error) {
	var final []validator.TestResult
	for i := range results {
		res, err := brunoToOpenAPI(results[i], router)
		if err != nil {
			return nil, err
		}
		final = append(final, res)
	}
	return final, nil
}

func brunoToOpenAPI(result Result, router routers.Router) (validator.TestResult, error) {
	request, err := translateRequest(result.Request, router)
	if err != nil {
		return validator.TestResult{}, err
	}
	response, err := translateResponse(result.Response, request.RequestValidationInput)
	if err != nil {
		return validator.TestResult{}, err
	}
	return validator.TestResult{
		Request:  request,
		Response: response,
		Id:       formatId(result.Test.Filename),
	}, nil
}

func formatId(filename string) string {
	filename = strings.TrimSuffix(filename, ".bru")
	filename = strings.TrimSuffix(filename, "-muted-")
	filename = strings.TrimSpace(filename)
	return filename
}

func translateRequest(brunoRequest Request, router routers.Router) (*validator.TestRequest, error) {
	// Translate request
	var requestBody io.Reader
	var prettyJSON bytes.Buffer
	if brunoRequest.Body != "" {
		requestBody = strings.NewReader(string(brunoRequest.Body))
		if err := json.Indent(&prettyJSON, []byte(brunoRequest.Body), "", "  "); err != nil {
			return nil, err
		}
	}

	httpReq, err := http.NewRequest(brunoRequest.Method, brunoRequest.Url, requestBody)
	if err != nil {
		return nil, err
	}
	// Set headers
	for header, value := range brunoRequest.Headers {
		httpReq.Header.Set(header, value)
	}
	route, pathParams, err := router.FindRoute(httpReq)
	if err != nil {
		return nil, err
	}
	// Disabling security checks
	route.Spec.Security = nil

	request := validator.TestRequest{
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request:    httpReq,
			PathParams: pathParams,
			Route:      route,
		},
		Body: prettyJSON.String(),
	}

	return &request, nil
}

func translateResponse(brunoResponse Response, request *openapi3filter.RequestValidationInput) (*validator.TestResponse, error) {
	headers := http.Header{}
	for header, value := range brunoResponse.Headers {
		headers.Set(header, value)
	}
	var bodyReader io.ReadCloser
	var prettyJSON bytes.Buffer
	if brunoResponse.Body != "" {
		bodyReader = io.NopCloser(strings.NewReader(string(brunoResponse.Body)))
		if err := json.Indent(&prettyJSON, []byte(brunoResponse.Body), "", "  "); err != nil {
			return nil, err
		}
	}

	return &validator.TestResponse{ResponseValidationInput: &openapi3filter.ResponseValidationInput{
		RequestValidationInput: request,
		Status:                 brunoResponse.Status,
		Header:                 headers,
		Body:                   bodyReader,
	},
		Body: prettyJSON.String(),
	}, nil
}
