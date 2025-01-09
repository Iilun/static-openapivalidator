package validator

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
)

const (
	Failure = "failure"
	Warning = "warning"
	Success = "success"
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
	requestResult, err := requestValidationResult(result, ctx)
	if err != nil {
		return nil, err
	}

	responseResult, err := responseValidationResult(result, ctx)
	if err != nil {
		return nil, err
	}

	return []ValidationResult{
		requestResult, responseResult,
	}, nil
}

func computeError(err error) (string, []ValidationError, error) {
	var multiError openapi3.MultiError
	if errors.As(err, &multiError) {
		return computeErrorFields(multiError)
	}
	return computeErrorFields([]error{err})
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
	status := Success
	var errorTitle string
	var err error
	var validationErrors []ValidationError
	if validationError != nil {
		status = Failure
		responseErr := new(openapi3filter.ResponseError)
		requestErr := new(openapi3filter.RequestError)
		if errors.As(validationError, &responseErr) {
			errorTitle = responseErr.Reason
			_, validationErrors, err = computeError(responseErr.Err)
			if err != nil {
				return "", "", nil, err
			}
		} else if errors.As(validationError, &requestErr) {
			errorTitle = requestErr.Reason
			_, validationErrors, err = computeError(requestErr.Err)
			if err != nil {
				return "", "", nil, err
			}
		} else {
			errorTitle, validationErrors, err = computeError(validationError)
			if err != nil {
				return "", "", nil, err
			}
		}
	}
	return status, errorTitle, validationErrors, nil
}

func requestValidationResult(result TestResult, ctx context.Context) (*RequestValidationResult, error) {
	var status, errAsString string
	var validationErrors []ValidationError
	var err error
	if result.Request.ParsingError != "" {
		status = Warning
		errAsString = result.Request.ParsingError
	} else {
		status, errAsString, validationErrors, err = computeResultFields(openapi3filter.ValidateRequest(ctx, result.Request.RequestValidationInput))
		if err != nil {
			return nil, err
		}
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

func responseValidationResult(result TestResult, ctx context.Context) (*ResponseValidationResult, error) {
	var status, errAsString string
	var validationErrors []ValidationError
	var err error

	if result.Response.ParsingError != "" {
		status = Warning
		errAsString = result.Response.ParsingError
	} else {
		status, errAsString, validationErrors, err = computeResultFields(openapi3filter.ValidateResponse(ctx, result.Response.ResponseValidationInput))
		if err != nil {
			return nil, err
		}
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
