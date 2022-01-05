package response

import (
	"encoding/json"

	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/codes"
)

type SuccessResponse struct {
	Code   codes.Code  `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type ErrResponse struct {
	Code   codes.Code  `json:"code"`
	Status string      `json:"status"`
	Err    interface{} `json:"error"`
}

func CreatedResponse(data interface{}) *httpbody.HttpBody {
	resp, _ := json.Marshal(SuccessResponse{
		Code:   codes.OK,
		Status: codes.OK.String(),
		Data:   data,
	})
	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        resp,
	}
}

func OKResponse(data interface{}) *httpbody.HttpBody {
	resp, _ := json.Marshal(SuccessResponse{
		Code:   codes.OK,
		Status: codes.OK.String(),
		Data:   data,
	})

	return &httpbody.HttpBody{
		ContentType: "application/x-protobuf",
		Data:        resp,
	}
}

func ErrNotFound(err interface{}) *httpbody.HttpBody {
	resp, _ := json.Marshal(ErrResponse{
		Code:   codes.NotFound,
		Status: codes.NotFound.String(),
		Err:    err,
	})

	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        resp,
	}
}

func ErrAlreadyExists(err interface{}) *httpbody.HttpBody {
	resp, _ := json.Marshal(ErrResponse{
		Code:   codes.AlreadyExists,
		Status: codes.AlreadyExists.String(),
		Err:    err,
	})

	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        resp,
	}
}
