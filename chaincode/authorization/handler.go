package authorization

import (
	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/handler"
	"github.com/lalloni/fabrikit/chaincode/response"
)

func Handler(action string, check Check, handler handler.Handler) handler.Handler {
	return func(ctx *context.Context) *response.Response {
		err := check(ctx)
		if err != nil {
			return response.Forbidden("%s forbidden: %s", action, err)
		}
		return handler(ctx)
	}
}
