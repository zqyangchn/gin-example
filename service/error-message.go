package service

import (
	"gin-example/pkg/errcode"
)

// Response struct
type ErrorMessageResponse struct {
	*errcode.ErrorMessage
	Data errcode.ErrorMessages
}

func GetAllErrorMessage() *ErrorMessageResponse {
	return &ErrorMessageResponse{
		ErrorMessage: errcode.Success,
		Data:         errcode.GetAllErrorMessage(),
	}
}
