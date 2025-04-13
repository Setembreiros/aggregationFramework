package api_connector

import "fmt"

type BadStatusCodeResponseError struct {
	statusCode int
}

func (e *BadStatusCodeResponseError) Error() string {
	errorMessage := fmt.Sprintf("Internal API call returned response with StatusCode: %d", e.statusCode)
	return errorMessage
}

func NewBadStatusCodeResponseError(statusCode int) *BadStatusCodeResponseError {
	return &BadStatusCodeResponseError{
		statusCode: statusCode,
	}
}

type ContentDeserializationError struct {
}

func (e *ContentDeserializationError) Error() string {
	errorMessage := fmt.Sprintf("Maping content failed")
	return errorMessage
}

func NewContentDeserializationError() *ContentDeserializationError {
	return &ContentDeserializationError{}
}
