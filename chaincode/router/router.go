package router

import (
	"fmt"
	"sort"

	"github.com/lalloni/fabrikit/chaincode/authorization"
	"github.com/lalloni/fabrikit/chaincode/handler"
)

type Name string

func (n *Name) String() string {
	return string(*n)
}

type Router interface {
	InitHandler() handler.Handler
	SetInitHandler(authorization.Check, handler.Handler)
	Handler(Name) handler.Handler
	SetHandler(Name, authorization.Check, handler.Handler)
	Functions() []Name
}

func New() Router {
	return &router{
		functionHandlers: map[string]handler.Handler{},
	}
}

type router struct {
	initHandler      handler.Handler
	functionHandlers map[string]handler.Handler
}

func (r *router) Functions() []Name {
	fs := []string(nil)
	for f := range r.functionHandlers {
		fs = append(fs, f)
	}
	sort.Strings(fs)
	ns := []Name(nil)
	for _, f := range fs {
		ns = append(ns, Name(f))
	}
	return ns
}

func (r *router) InitHandler() handler.Handler {
	return r.initHandler
}

func (r *router) SetInitHandler(ch authorization.Check, h handler.Handler) {
	if h == nil {
		h = handler.SuccessHandler
	}
	if ch != nil {
		r.initHandler = authorization.Handler("init", ch, h)
	} else {
		r.initHandler = h
	}
}

func (r *router) Handler(n Name) handler.Handler {
	return r.functionHandlers[n.String()]
}

func (r *router) SetHandler(n Name, ch authorization.Check, h handler.Handler) {
	if h == nil {
		h = handler.SuccessHandler
	}
	if ch != nil {
		r.functionHandlers[n.String()] = authorization.Handler(fmt.Sprintf("invoke function %q", n), ch, h)
	} else {
		r.functionHandlers[n.String()] = h
	}
}
