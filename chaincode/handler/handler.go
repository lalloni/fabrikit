package handler

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/response"
)

type Handler func(*context.Context) *response.Response

func Name(h Handler) string {
	s := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	s = s[strings.LastIndex(s, ".")+1:]
	s = strings.TrimSuffix(s, "Handler")
	return s
}
