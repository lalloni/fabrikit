package param

import (
	"reflect"
	"strconv"
	"unicode/utf8"

	"github.com/pkg/errors"
)

func Uint64Var(ref *uint64) TypedParam {
	return Typed("natural integer",
		reflect.TypeOf(uint64(0)),
		func(arg []byte) (interface{}, error) {
			r, err := parseuint64(arg)
			if err != nil {
				return nil, err
			}
			if ref != nil {
				*ref = r
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

func StringVar(ref *string) TypedParam {
	return Typed("string", reflect.TypeOf(""), func(bs []byte) (interface{}, error) {
		if !utf8.Valid(bs) {
			return nil, errors.New("invalid UTF-8 string")
		}
		s := string(bs)
		if ref != nil {
			*ref = s
		}
		return s, nil
	})
}

var String = StringVar(nil)

func BoolVar(ref *bool) TypedParam {
	return Typed("boolean", reflect.TypeOf(true), func(bs []byte) (interface{}, error) {
		s := string(bs)
		var v bool
		switch s {
		case "true":
			v = true
		case "false":
			v = false
		default:
			return nil, errors.New("invalid boolean")
		}
		if ref != nil {
			*ref = v
		}
		return v, nil
	})
}

var Bool = BoolVar(nil)
