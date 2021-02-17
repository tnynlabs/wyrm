package utils

import (
	"fmt"
)

type ServiceErrCode string

const UnexpectedCode = ServiceErrCode("UNEXPECTED")

type ServiceErr struct {
	Code    ServiceErrCode `json:"code"`
	Message string         `json:"message"`
}

func (e *ServiceErr) Error() string {
	return fmt.Sprintf("error: [%s] %s", e.Code, e.Message)
}

func ToServiceErr(err error) *ServiceErr {
	if err == nil {
		return nil
	}

	v, ok := err.(*ServiceErr)
	if ok {
		return v
	}

	return &ServiceErr{
		Code:    UnexpectedCode,
		Message: "An unexpected error occurred",
	}
}
