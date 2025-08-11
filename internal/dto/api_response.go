package dto

import (
	"github.com/frenswifbenefits/myfren/internal/openapi"
)

const (
	SuccessCode = 0
	ErrorCode   = 1

	SuccessMessage = "Success"
)

func MakeSuccessAPIResponse() api_types.Response {
	return api_types.Response{
		Code:    SuccessCode,
		Message: SuccessMessage,
	}
}

func MakeErrorAPIResponse(err error) api_types.Response {
	return api_types.Response{
		Code:    ErrorCode,
		Message: err.Error(),
	}
}

func MakeBadRequestAPIResponse(err error) api_types.Response {
	return api_types.Response{
		Code:    ErrorCode,
		Message: err.Error(),
	}
}

func MakeForbiddenAPIResponse(err error) api_types.Response {
	return api_types.Response{
		Code:    ErrorCode,
		Message: err.Error(),
	}
}

func MakeNotFoundAPIResponse(err error) api_types.Response {
	return api_types.Response{
		Code:    ErrorCode,
		Message: err.Error(),
	}
}

func MakeBadApiKeyAPIResponse(err error) api_types.Response {
	return api_types.Response{
		Code:    10003,
		Message: err.Error(),
	}
}
