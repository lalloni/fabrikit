package router

import (
	"github.com/lalloni/fabrikit/chaincode/authorization"
	"github.com/lalloni/fabrikit/chaincode/handler"
)

type Route struct {
	Name    string
	Check   authorization.Check
	Handler handler.Handler
}

type Config struct {
	Init *Route
	Funs []Route
}
