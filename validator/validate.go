package validator

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"static-openapivalidator/logger"
)

const (
	Failure = "failure"
	Warning = "warning"
	Ignored = "ignored"
	Success = "success"
)

func Validate(results []TestResult, ctx context.Context) ([]ValidationResult, error) {
	var final []ValidationResult
	logger.Log("Validator: Validating %d results", len(results))
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
	var results []ValidationResult

	requestResult, err := requestValidationResult(result, ctx)
	if err != nil {
		return nil, errors.New("error validating request: " + err.Error())
	}
	if requestResult != nil {
		results = append(results, requestResult)
	}

	responseResult, err := responseValidationResult(result, ctx)
	if err != nil {
		return nil, errors.New("error validating response: " + err.Error())
	}
	if responseResult != nil {
		results = append(results, responseResult)
	}
	return results, nil
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
			if errAsString == "" && multiError[i] != nil {
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
			if responseErr.Reason != "" {
				errorTitle = responseErr.Reason
			} else {
				errorTitle = validationError.Error()
			}
			_, validationErrors, err = computeError(responseErr.Err)
			if err != nil {
				return "", "", nil, err
			}
		} else if errors.As(validationError, &requestErr) {
			if requestErr.Reason != "" {
				errorTitle = requestErr.Reason
			} else {
				errorTitle = validationError.Error()
			}
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
		if len(validationErrors) == 0 {
			errorTitle = validationError.Error()
		}
	}
	return status, errorTitle, validationErrors, nil
}

func requestValidationResult(result TestResult, ctx context.Context) (*RequestValidationResult, error) {
	var status, errAsString string
	var validationErrors []ValidationError
	var err error

	if result.Request.Ignored {
		status = Ignored
	} else {
		if result.Request.ParsingError != "" {
			status = Warning
			errAsString = result.Request.ParsingError
		} else {
			status, errAsString, validationErrors, err = computeResultFields(openapi3filter.ValidateRequest(ctx, result.Request.RequestValidationInput))
			if err != nil {
				return nil, err
			}
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
		Method:       result.Request.Request.Method,
	}, nil
}

func responseValidationResult(result TestResult, ctx context.Context) (*ResponseValidationResult, error) {
	var status, errAsString string
	var validationErrors []ValidationError
	var err error

	if result.Response.Ignored {
		status = Ignored
	} else {
		if result.Response.ParsingError != "" {
			status = Warning
			errAsString = result.Response.ParsingError
		} else {
			status, errAsString, validationErrors, err = computeResultFields(openapi3filter.ValidateResponse(ctx, result.Response.ResponseValidationInput))
			if err != nil {
				return nil, err
			}
		}
	}

	return &ResponseValidationResult{
		TestId:       result.Id,
		Url:          result.Request.Request.URL.Path,
		ErrorSummary: errAsString,
		Errors:       validationErrors,
		Status:       status,
		Code:         result.Response.Status,
		Body:         result.Response.Body,
		Headers:      result.Response.ResponseValidationInput.Header,
	}, nil
}
