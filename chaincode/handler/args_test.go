package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lalloni/fabrikit/chaincode/handler/param"
)

func TestExtract(t *testing.T) {
	type args struct {
		args [][]byte
		pars []param.Param
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{"no pars", args{bytes(), pars()}, wants(), false},
		{"missing 1 arg", args{bytes(), pars(param.Uint64)}, wants(), true},
		{"missing 2 args", args{bytes(), pars(param.Uint64, param.String)}, wants(), true},
		{"missing 1 args", args{bytes("1"), pars(param.Uint64, param.String)}, wants(), true},
		{"get uint64, string", args{bytes("1", "bla"), pars(param.Uint64, param.String)}, wants(uint64(1), "bla"), false},
		{"1 extra arg", args{bytes("1"), pars()}, wants(), true},
		{"2 extra args", args{bytes("1", "2"), pars()}, wants(), true},
		{"1 extra args", args{bytes("1", "2"), pars(param.Uint64)}, wants(), true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			got, err := ExtractArgs(tt.args.args, tt.args.pars...)
			if tt.wantErr {
				a.Error(err)
			} else {
				a.NoError(err)
			}
			a.EqualValues(tt.want, got)
		})
	}
}

func wants(vv ...interface{}) []interface{} {
	return vv
}

func bytes(ss ...string) [][]byte {
	bs := [][]byte(nil)
	for _, s := range ss {
		bs = append(bs, []byte(s))
	}
	return bs
}

func pars(ps ...param.Param) []param.Param {
	return ps
}
