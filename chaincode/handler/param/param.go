package param

import "reflect"

type Parser func(arg []byte) (value interface{}, err error)

type Transformer func(in interface{}) (out interface{}, err error)

type Param interface {
	Name() string
	From(arg []byte) (value interface{}, err error)
}

func New(name string, p Parser) Param {
	return &param{name: name, parser: p}
}

type param struct {
	name   string
	parser Parser
}

func (p *param) Name() string {
	return p.name
}

func (p *param) From(arg []byte) (interface{}, error) {
	return p.parser(arg)
}

type TypedParam interface {
	Param
	Type() reflect.Type
}

func Typed(name string, t reflect.Type, p Parser) TypedParam {
	return &typed{param: param{name: name, parser: p}, typ: t}
}

type typed struct {
	param
	typ reflect.Type
}

func (t *typed) Type() reflect.Type {
	return t.typ
}

func Untyped(tt ...TypedParam) []Param {
	pp := []Param(nil)
	for _, t := range tt {
		pp = append(pp, t)
	}
	return pp
}

func Specialize(p Param, name string, next Transformer) Param {
	return New(name, combine(p.From, next))
}

func SpecializeTyped(p TypedParam, name string, next Transformer) TypedParam {
	return Typed(name, p.Type(), combine(p.From, next))
}

func combine(p Parser, t Transformer) Parser {
	return func(arg []byte) (interface{}, error) {
		r, err := p(arg)
		if err != nil {
			return nil, err
		}
		r, err = t(r)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
}
