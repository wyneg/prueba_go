package models

import (
	"fmt"
	"net/http"
	"strconv"
)

type CustomError struct {
	Code     string `json:"code"`
	ErrorMsg string `json:"error"`
	Err      error  `json:"-"`
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.ErrorMsg)
}

func NewNotFoundError(message string) *CustomError {
	return &CustomError{
		Code:     strconv.Itoa(http.StatusNotFound),
		ErrorMsg: message,
	}
}

func NewBadRequestError(message string) *CustomError {
	return &CustomError{
		Code:     strconv.Itoa(http.StatusBadRequest),
		ErrorMsg: message,
	}
}
