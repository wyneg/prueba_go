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

func NewError(httpError int, message string) *CustomError {
	return &CustomError{
		Code:     strconv.Itoa(httpError),
		ErrorMsg: message,
	}
}

func NewNotFoundError(message string) *CustomError {
	return NewError(http.StatusNotFound, message)
}

func NewBadRequestError(message string) *CustomError {
	return NewError(http.StatusBadRequest, message)
}

func NewInternalServerError(message string) *CustomError {
	return NewError(http.StatusInternalServerError, message)
}

func NewConflictError(message string) *CustomError {
	return NewError(http.StatusConflict, message)
}
