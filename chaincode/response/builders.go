package response

import (
	"fmt"

	"github.com/lalloni/fabrikit/chaincode/response/status"
)

func OK(result interface{}) *Response {
	return StatusWithResult(status.OK, result)
}

func NotFound() *Response {
	return Status(status.NotFound)
}

func NotFoundWithMessage(msg string, args ...interface{}) *Response {
	return StatusWithMessage(status.NotFound, msg, args...)
}

func Forbidden(msg string, args ...interface{}) *Response {
	return StatusWithMessage(status.Forbidden, msg, args...)
}

func BadRequest(msg string, args ...interface{}) *Response {
	return StatusWithMessage(status.BadRequest, msg, args...)
}

func BadRequestWithFault(fault interface{}) *Response {
	return StatusWithFault(status.BadRequest, fault)
}

func Error(msg string, args ...interface{}) *Response {
	return StatusWithMessage(status.Error, msg, args...)
}

func NotImplemented(function string) *Response {
	return StatusWithMessage(status.BadRequest, "function '%s' is not implemented", function)
}

func StatusWithResult(status int32, result interface{}) *Response {
	return &Response{
		Status: status,
		Payload: &Payload{
			Content: result,
		},
	}
}

func StatusWithFault(status int32, fault interface{}) *Response {
	return &Response{
		Status: status,
		Payload: &Payload{
			Fault: fault,
		},
	}
}

func StatusWithMessage(status int32, msg string, args ...interface{}) *Response {
	return &Response{
		Status:  status,
		Message: fmt.Sprintf(msg, args...),
	}
}

func Status(status int32) *Response {
	return &Response{
		Status: status,
	}
}
