package handler

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/handler/param"
)

func TestFunc(t *testing.T) {
	type args struct {
		function interface{}
		pars     []param.TypedParam
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"not a function", args{"bla", nil}, true},
		{"bad return 1", args{func() (string, int, int) { return "", 0, 0 }, nil}, true},
		{"bad return 2", args{func() {}, nil}, true},
		{"return ok", args{func(*context.Context) int { return 0 }, nil}, false},
		{"no ctx", args{func(int) int { return 0 }, nil}, true},
		{"bad args from params 1", args{func(*context.Context, int, string) int { return 0 }, nil}, true},
		{"1 param ok", args{func(*context.Context, int) int { return 0 }, []param.TypedParam{param.Typed("integer", reflect.TypeOf(0), nil)}}, false},
		{"bad params 1", args{func(*context.Context, string) int { return 0 }, []param.TypedParam{param.Typed("integer", reflect.TypeOf(0), nil)}}, true},
		{"bad params 2", args{func(*context.Context, string, int) int { return 0 }, []param.TypedParam{param.Typed("integer", reflect.TypeOf(0), nil)}}, true},
		{"2 params ok", args{func(*context.Context, string, int) int { return 0 }, []param.TypedParam{param.Typed("string", reflect.TypeOf(""), nil), param.Typed("integer", reflect.TypeOf(0), nil)}}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			v, err := Func(tt.args.function, tt.args.pars...)
			if tt.wantErr {
				a.Error(err)
				a.Nil(v)
			} else {
				a.NoError(err)
				a.NotNil(v)
			}
		})
	}
}
