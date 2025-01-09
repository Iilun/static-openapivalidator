package validator

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
)

func Validate(results []TestResult, ctx context.Context) ([]ValidationResult, error) {
	var final []ValidationResult
	for i := range results {
		validationResults, err := validateResult(results[i], ctx)
		if err != nil {
			return nil, err
		}
		final = append(final, validationResults...)
	}
	return final, nil
}

// ctx contains the openapi spec
func validateResult(result TestResult, ctx context.Context) ([]ValidationResult, error) {
	requestResult, err := requestValidationResult(result, openapi3filter.ValidateRequest(ctx, result.Request.RequestValidationInput))
	if err != nil {
		return nil, err
	}

	responseResult, err := responseValidationResult(result, openapi3filter.ValidateResponse(ctx, result.Response.ResponseValidationInput))
	if err != nil {
		return nil, err
	}

	return []ValidationResult{
		requestResult, responseResult,
	}, nil
}

func computeErrorFields(multiError []error) (string, []ValidationError, error) {
	var errAsString string
	var validationErrors []ValidationError
	for i := range multiError {

		var schemaError *openapi3.SchemaError
		if errors.As(multiError[i], &schemaError) {
			if errAsString != "" && errAsString != schemaError.Reason {
				errAsString = "Multiple errors"
			} else {
				errAsString = schemaError.Reason
			}
			schemaBytes, err := json.MarshalIndent(schemaError.Schema, "", "  ")
			if err != nil {
				return "", nil, err
			}
			validationErrors = append(validationErrors, ValidationError{
				Title:  schemaError.Error(),
				Schema: string(schemaBytes),
			})
		} else {
			// Do not erase a more important error
			if errAsString == "" {
				errAsString = multiError[i].Error()
			}
		}
	}
	if len(validationErrors) == 1 {
		errAsString = validationErrors[0].Title
	}

	return errAsString, validationErrors, nil
}

func computeResultFields(validationError error) (string, string, []ValidationError, error) {
	status := "success"
	var errorTitle string
	var err error
	var validationErrors []ValidationError
	if validationError != nil {
		status = "failure"
		var multiError openapi3.MultiError
		if errors.As(validationError, &multiError) {
			errorTitle, validationErrors, err = computeErrorFields(multiError)
			if err != nil {
				return "", "", nil, err
			}
		} else {
			errorTitle, validationErrors, err = computeErrorFields([]error{validationError})
			if err != nil {
				return "", "", nil, err
			}
		}
	}
	return status, errorTitle, validationErrors, nil
}

func requestValidationResult(result TestResult, validationError error) (*RequestValidationResult, error) {
	status, errAsString, validationErrors, err := computeResultFields(validationError)
	if err != nil {
		return nil, err
	}

	return &RequestValidationResult{
		TestId:       result.Id,
		Url:          result.Request.Request.URL.Path,
		ErrorSummary: errAsString,
		Errors:       validationErrors,
		Status:       status,
		Body:         result.Request.Body,
		Headers:      result.Request.Request.Header,
	}, nil
}

func responseValidationResult(result TestResult, validationError error) (*ResponseValidationResult, error) {
	status, errAsString, validationErrors, err := computeResultFields(validationError)
	if err != nil {
		return nil, err
	}

	return &ResponseValidationResult{
		TestId:       result.Id,
		Url:          result.Request.Request.URL.Path,
		ErrorSummary: errAsString,
		Errors:       validationErrors,
		Status:       status,
		Body:         result.Response.Body,
		Headers:      result.Response.ResponseValidationInput.Header,
	}, nil
}