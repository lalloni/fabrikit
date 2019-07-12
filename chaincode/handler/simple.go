package handler

import (
	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/response"
)

func SuccessHandler(ctx *context.Context) *response.Response {
	return response.OK(nil)
}

func EchoHandler(ctx *context.Context) *response.Response {
	data, err := ctx.ArgBytes(1)
	if err != nil {
		return response.BadRequest(err.Error())
	}
	return response.OK(data)
}

func ValueHandler(v interface{}) Handler {
	return func(ctx *context.Context) *response.Response {
		if err := CheckArgsCount(ctx, 0); err != nil {
			return response.BadRequest(err.Error())
		}
		return response.OK(v)
	}
}

func NotImplementedHandler(ctx *context.Context) *response.Response {
	return response.NotImplemented(ctx.Function())
}
