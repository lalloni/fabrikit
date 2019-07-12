package param

import (
	"reflect"
	"strconv"
	"unicode/utf8"

	"github.com/pkg/errors"
)

func Uint64Var(v *uint64) TypedParam {
	return Typed("natural integer",
		reflect.TypeOf(uint64(0)),
		func(arg []byte) (interface{}, error) {
			r, err := parseuint64(arg)
			if err != nil {
				return nil, err
			}
			if v != nil {
				*v = r
			}
			return r, err
		})
}

var Uint64 = Uint64Var(nil)

func parseuint64(arg []byte) (uint64, error) {
	s := string(arg)
	r, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		if e, ok := err.(*strconv.NumError); ok {
			err = e.Err
		}
		return 0, errors.Errorf("invalid natural integer: %v: '%v'", err, s)
	}
	return r, nil
}

var String = Typed("string", reflect.TypeOf(""), parsestring)

func parsestring(arg []byte) (interface{}, error) {
	if !utf8.Valid(arg) {
		return nil, errors.Errorf("invalid UTF-8 string")
	}
	return string(arg), nil
}
