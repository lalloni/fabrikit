package router

import (
	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/handler"
	"github.com/lalloni/fabrikit/chaincode/response"
)

func FunctionsHandler(r Router) handler.Handler {
	return func(ctx *context.Context) *response.Response {
		if err := handler.CheckArgsCount(ctx, 0); err != nil {
			return response.BadRequest(err.Error())
		}
		return response.OK(r.Functions())
	}
}
