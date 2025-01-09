package validator

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3filter"
)

type TestResult struct {
	AdditionalInfos map[string]string
	Request         *TestRequest
	Response        *TestResponse
	Id              string
}

type TestRequest struct {
	*openapi3filter.RequestValidationInput
	Body string
}

type TestResponse struct {
	*openapi3filter.ResponseValidationInput
	Body string
}

type ValidationResult interface {
	GetType() string
	GetTestId() string
	GetErrorSummary() string
	GetErrors() []ValidationError
	GetUrl() string
	GetBody() string
	GetHeaders() map[string][]string
}

type ValidationError struct {
	Title  string `json:"title,omitempty"`
	Schema string `json:"schema,omitempty"`
}

type RequestValidationResult struct {
	TestId       string
	Url          string
	ErrorSummary string
	Errors       []ValidationError
	Status       string
	Body         string
	Headers      map[string][]string
}

func (r RequestValidationResult) GetType() string {
	return "request"
}

func (r RequestValidationResult) GetUrl() string {
	return r.Url
}

func (r RequestValidationResult) GetTestId() string {
	return r.TestId
}

func (r RequestValidationResult) GetErrorSummary() string {
	return r.ErrorSummary
}

func (r RequestValidationResult) GetErrors() []ValidationError {
	return r.Errors
}

func (r RequestValidationResult) GetBody() string {
	return r.Body
}

func (r RequestValidationResult) GetHeaders() map[string][]string {
	return r.Headers
}

func (r RequestValidationResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonValidationResult{
		TestId:       r.TestId,
		Url:          r.Url,
		ErrorSummary: r.ErrorSummary,
		Errors:       r.Errors,
		Status:       r.Status,
		Type:         r.GetType(),
		Body:         r.Body,
		Headers:      r.Headers,
	})
}

type ResponseValidationResult struct {
	TestId       string
	Url          string
	ErrorSummary string
	Errors       []ValidationError
	Status       string
	Body         string
	Headers      map[string][]string
}

func (r ResponseValidationResult) GetType() string {
	return "response"
}

func (r ResponseValidationResult) GetUrl() string {
	return r.Url
}

func (r ResponseValidationResult) GetTestId() string {
	return r.TestId
}

func (r ResponseValidationResult) GetErrorSummary() string {
	return r.ErrorSummary
}

func (r ResponseValidationResult) GetErrors() []ValidationError {
	return r.Errors
}

func (r ResponseValidationResult) GetBody() string {
	return r.Body
}

func (r ResponseValidationResult) GetHeaders() map[string][]string {
	return r.Headers
}

func (r ResponseValidationResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonValidationResult{
		TestId:       r.TestId,
		Url:          r.Url,
		ErrorSummary: r.ErrorSummary,
		Errors:       r.Errors,
		Status:       r.Status,
		Type:         r.GetType(),
		Body:         r.Body,
		Headers:      r.Headers,
	})
}

type jsonValidationResult struct {
	TestId       string              `json:"id"`
	Url          string              `json:"url"`
	ErrorSummary string              `json:"error,omitempty"`
	Errors       []ValidationError   `json:"errors,omitempty"`
	Status       string              `json:"status"`
	Type         string              `json:"type"`
	Body         string              `json:"body"`
	Headers      map[string][]string `json:"headers"`
}