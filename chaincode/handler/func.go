package handler

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/pkg/errors"

	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/handler/param"
	"github.com/lalloni/fabrikit/chaincode/response"
)

var contextType = reflect.TypeOf(&context.Context{})

func MustFunc(function interface{}, pars ...param.TypedParam) Handler {
	h, err := Func(function, pars...)
	if err != nil {
		panic(err)
	}
	return h
}

func Func(function interface{}, pars ...param.TypedParam) (Handler, error) {

	fun := reflect.ValueOf(function)
	funType := fun.Type()

	if fun.Kind() != reflect.Func {
		return nil, errors.Errorf("not a function: %v", function)
	}

	funName := runtime.FuncForPC(fun.Pointer()).Name()

	if funType.NumOut() != 1 {
		return nil, errors.Errorf("function %s must return 1 value", funName)
	}

	cardinality := len(pars)

	if cardinality+1 != funType.NumIn() {
		s := ""
		if cardinality+1 != 1 {
			s = "s"
		}
		return nil, errors.Errorf("function %s must have %d parameter%s of type %s", funName, cardinality+1, s, types(pars))
	}

	if funType.In(0) != contextType {
		return nil, errors.Errorf("function %s parameter 0 must be a %s", funName, contextType)
	}

	for i, par := range pars {
		t := funType.In(i + 1)
		if !t.AssignableTo(par.Type()) {
			return nil, errors.Errorf("function %s parameter %d must be %s (assignable to %s type) but is %s", funName, i+1, par.Name(), par.Type(), t)
		}
	}

	return func(ctx *context.Context) *response.Response {
		args, err := ExtractArgs(ctx.Stub.GetArgs()[1:], param.Untyped(pars...)...)
		if err != nil {
			return response.BadRequest(err.Error())
		}
		vals := []reflect.Value{reflect.ValueOf(ctx)}
		for _, arg := range args {
			vals = append(vals, reflect.ValueOf(arg))
		}
		ret := fun.Call(vals)[0].Interface()
		if res, ok := ret.(*response.Response); ok {
			return res
		}
		return response.OK(ret)
	}, nil

}

func types(pars []param.TypedParam) string {
	ss := []string{contextType.String()}
	for _, par := range pars {
		ss = append(ss, par.Type().String())
	}
	return strings.Join(ss, ", ")
}
