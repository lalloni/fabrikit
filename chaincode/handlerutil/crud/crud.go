package crud

import (
	"strings"

	auth "github.com/lalloni/fabrikit/chaincode/authorization"
	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/handler"
	"github.com/lalloni/fabrikit/chaincode/handler/param"
	"github.com/lalloni/fabrikit/chaincode/response"
	"github.com/lalloni/fabrikit/chaincode/router"
	"github.com/lalloni/fabrikit/chaincode/store"
)

type Validator func(*context.Context, interface{}) *response.Response

type opt struct {
	get      bool
	getall   bool
	getrange bool
	has      bool
	put      bool
	putlist  bool
	del      bool
	delrange bool

	getsingleton      bool
	getcollection     bool
	getcollectionitem bool

	defaultcheck auth.Check
	readcheck    auth.Check
	writecheck   auth.Check

	getcheck      auth.Check
	getallcheck   auth.Check
	getrangecheck auth.Check
	hascheck      auth.Check
	putcheck      auth.Check
	putlistcheck  auth.Check
	delcheck      auth.Check
	delrangecheck auth.Check

	getsingletoncheck      auth.Check
	getcollectioncheck     auth.Check
	getcollectionitemcheck auth.Check

	putvalidator Validator

	id   param.Param
	item param.Param
	list param.Param
}

type Option func(*opt)

func WithGet(b bool) Option      { return func(o *opt) { o.get = b } }
func WithGetAll(b bool) Option   { return func(o *opt) { o.getall = b } }
func WithGetRange(b bool) Option { return func(o *opt) { o.getrange = b } }
func WithHas(b bool) Option      { return func(o *opt) { o.has = b } }
func WithPut(b bool) Option      { return func(o *opt) { o.put = b } }
func WithPutList(b bool) Option  { return func(o *opt) { o.putlist = b } }
func WithDel(b bool) Option      { return func(o *opt) { o.del = b } }
func WithDelRange(b bool) Option { return func(o *opt) { o.delrange = b } }

func WithGetSingleton(b bool) Option      { return func(o *opt) { o.getsingleton = b } }
func WithGetCollection(b bool) Option     { return func(o *opt) { o.getcollection = b } }
func WithGetCollectionItem(b bool) Option { return func(o *opt) { o.getcollectionitem = b } }

func WithDefaultCheck(c auth.Check) Option { return func(o *opt) { o.defaultcheck = c } }
func WithReadCheck(c auth.Check) Option    { return func(o *opt) { o.readcheck = c } }
func WithWriteCheck(c auth.Check) Option   { return func(o *opt) { o.writecheck = c } }

func WithGetCheck(c auth.Check) Option      { return func(o *opt) { o.getcheck = c } }
func WithGetAllCheck(c auth.Check) Option   { return func(o *opt) { o.getallcheck = c } }
func WithGetRangeCheck(c auth.Check) Option { return func(o *opt) { o.getrangecheck = c } }
func WithHasCheck(c auth.Check) Option      { return func(o *opt) { o.hascheck = c } }
func WithPutCheck(c auth.Check) Option      { return func(o *opt) { o.putcheck = c } }
func WithPutListCheck(c auth.Check) Option  { return func(o *opt) { o.putlistcheck = c } }
func WithDelCheck(c auth.Check) Option      { return func(o *opt) { o.delcheck = c } }
func WithDelRangeCheck(c auth.Check) Option { return func(o *opt) { o.delrangecheck = c } }

func WithGetSingletonCheck(c auth.Check) Option  { return func(o *opt) { o.getsingletoncheck = c } }
func WithGetCollectionCheck(c auth.Check) Option { return func(o *opt) { o.getcollectioncheck = c } }
func WithGetCollectionItemCheck(c auth.Check) Option {
	return func(o *opt) { o.getcollectionitemcheck = c }
}

func WithPutValidator(v Validator) Option { return func(o *opt) { o.putvalidator = v } }

func WithIDParam(p param.Param) Option   { return func(o *opt) { o.id = p } }
func WithItemParam(p param.Param) Option { return func(o *opt) { o.item = p } }
func WithListParam(p param.Param) Option { return func(o *opt) { o.list = p } }

func WithDefaults() Option {
	return func(o *opt) {
		for _, option := range Defaults {
			option(o)
		}
	}
}

var Defaults = []Option{
	WithHas(true),
	WithGet(true),
	WithGetAll(true),
	WithGetRange(true),
	WithPut(true),
	WithPutList(true),
	WithGetSingleton(true),
	WithGetCollection(true),
	WithGetCollectionItem(true),
	WithDel(true),
	WithDelRange(true),
	WithDefaultCheck(auth.Allowed),
}

func AddHandlers(r router.Router, s *store.Schema, opts ...Option) {
	o := &opt{}
	for _, opt := range opts {
		opt(o)
	}
	name := strings.Title(s.Name())
	if o.get {
		c := pri(o.getcheck, o.readcheck, o.defaultcheck)
		add(r, "Get"+name, c, GetHandler(s, o.id))
	}
	if o.getall {
		c := pri(o.getallcheck, o.readcheck, o.defaultcheck)
		add(r, "Get"+name+"All", c, GetAllHandler(s))
	}
	if o.getrange {
		c := pri(o.getrangecheck, o.readcheck, o.defaultcheck)
		add(r, "Get"+name+"Range", c, GetRangeHandler(s, o.id))
	}
	if o.has {
		c := pri(o.hascheck, o.readcheck, o.defaultcheck)
		add(r, "Has"+name, c, HasHandler(s, o.id))
	}
	if o.put {
		c := pri(o.putcheck, o.writecheck, o.defaultcheck)
		add(r, "Put"+name, c, PutHandler(s, o.item, o.putvalidator))
	}
	if o.putlist {
		c := pri(o.putlistcheck, o.writecheck, o.defaultcheck)
		add(r, "Put"+name+"List", c, PutListHandler(s, o.list, o.putvalidator))
	}
	if o.del {
		c := pri(o.delcheck, o.writecheck, o.defaultcheck)
		add(r, "Del"+name, c, DelHandler(s, o.id))
	}
	if o.delrange {
		c := pri(o.delcheck, o.writecheck, o.defaultcheck)
		add(r, "Del"+name+"Range", c, DelRangeHandler(s, o.id))
	}
	if o.getsingleton {
		c := pri(o.getsingletoncheck, o.readcheck, o.defaultcheck)
		for _, singleton := range s.Singletons() {
			add(r, "Get"+name+strings.Title(singleton.Name), c, GetSingletonHandler(s, singleton, o.id))
		}
	}
	if o.getcollection {
		c := pri(o.getcollectioncheck, o.readcheck, o.defaultcheck)
		for _, collection := range s.Collections() {
			add(r, "Get"+name+strings.Title(collection.Name), c, GetCollectionHandler(s, collection, o.id))
		}
	}
	if o.getcollectionitem {
		c := pri(o.getcollectioncheck, o.readcheck, o.defaultcheck)
		for _, collection := range s.Collections() {
			add(r, "Get"+name+strings.Title(collection.Name)+"Item", c, GetCollectionItemHandler(s, collection, o.id))
		}
	}
}

func pri(cs ...auth.Check) auth.Check {
	for _, c := range cs {
		if c != nil {
			return c
		}
	}
	return auth.Forbidden
}

func add(r router.Router, name string, c auth.Check, h handler.Handler) {
	r.SetHandler(router.Name(name), c, h)
}

func GetHandler(s *store.Schema, id param.Param) handler.Handler {
	return func(c *context.Context) *response.Response {
		args, err := handler.ExtractArgs(c.Args(), id)
		if err != nil {
			return response.BadRequest("invalid %s id: %v", s.Name(), err)
		}
		v, err := c.Store.GetComposite(s, args[0])
		if err != nil {
			return response.Error("getting %s: %v", s.Name(), err)
		}
		if v == nil {
			return response.NotFoundWithMessage("%s identified with %v not found", s.Name(), args[0])
		}
		return response.OK(v)
	}
}

func GetAllHandler(s *store.Schema) handler.Handler {
	return func(c *context.Context) *response.Response {
		_, err := handler.ExtractArgs(c.Args()) // no parameters
		if err != nil {
			return response.BadRequest(err.Error())
		}
		v, err := c.Store.GetCompositeAll(s)
		if err != nil {
			return response.Error("getting %s: %v", s.Name(), err)
		}
		return response.OK(v)
	}
}

func GetRangeHandler(s *store.Schema, id param.Param) handler.Handler {
	return func(c *context.Context) *response.Response {
		args, err := handler.ExtractArgs(c.Args(), id, id)
		if err != nil {
			return response.BadRequest("invalid %s id: %v", s.Name(), err)
		}
		v, err := c.Store.GetCompositeRange(s, store.R(args[0], args[1]))
		if err != nil {
			return response.Error("getting %s range: %v", s.Name(), err)
		}
		return response.OK(v)
	}
}

func PutHandler(s *store.Schema, val param.Param, valid Validator) handler.Handler {
	return func(c *context.Context) *response.Response {
		args, err := handler.ExtractArgs(c.Args(), val)
		if err != nil {
			return response.BadRequest("invalid %s: %v", s.Name(), err)
		}
		if valid != nil {
			res := valid(c, args[0])
			if res != nil {
				return res
			}
		}
		err = c.Store.PutComposite(s, args[0])
		if err != nil {
			return response.Error("putting %s: %v", s.Name(), err)
		}
		return response.OK(nil)
	}
}

func DelHandler(s *store.Schema, id param.Param) handler.Handler {
	return func(c *context.Context) *response.Response {
		args, err := handler.ExtractArgs(c.Args(), id)
		if err != nil {
			return response.BadRequest("invalid %s id: %v", s.Name(), err)
		}
		exist, err := c.Store.HasComposite(s, args[0])
		if err != nil {
			return response.Error("checking %s existence: %v", s.Name(), err)
		}
		if !exist {
			return response.NotFoundWithMessage("%s identified with %v not found", s.Name(), args[0])
		}
		err = c.Store.DelComposite(s, args[0])
		if err != nil {
			return response.Error("deleting %s: %v", s.Name(), err)
		}
		return response.OK(nil)
	}
}

func DelRangeHandler(s *store.Schema, id param.Param) handler.Handler {
	return func(c *context.Context) *response.Response {
		args, err := handler.ExtractArgs(c.Args(), id, id)
		if err != nil {
			return response.BadRequest("invalid %s id: %v", s.Name(), err)
		}
		v, err := c.Store.DelCompositeRange(s, store.R(args[0], args[1]))
		if err != nil {
			return response.Error("getting %s range: %v", s.Name(), err)
		}
		return response.OK(v)
	}
}

func HasHandler(s *store.Schema, id param.Param) handler.Handler {
	return func(c *context.Context) *response.Response {
		args, err := handler.ExtractArgs(c.Args(), id)
		if err != nil {
			return response.BadRequest("invalid %s id: %v", s.Name(), err)
		}
		exist, err := c.Store.HasComposite(s, args[0])
		if err != nil {
			return response.Error("getting %s existence: %v", s.Name(), err)
		}
		return response.OK(exist)
	}
}

func PutListHandler(s *store.Schema, list param.Param, valid Validator) handler.Handler {
	return func(c *context.Context) *response.Response {
		args, err := handler.ExtractArgs(c.Args(), list)
		if err != nil {
			return response.BadRequest("invalid %s list: %v", s.Name(), err)
		}
		count := 0
		for _, v := range args[0].([]interface{}) {
			if valid != nil {
				res := valid(c, v)
				if res != nil {
					return res
				}
			}
			err := c.Store.PutComposite(s, v)
			if err != nil {
				return response.Error("putting %s: %v", s.Name(), err)
			}
			count++
		}
		return response.OK(count)
	}
}

func GetSingletonHandler(s *store.Schema, singleton *store.Singleton, id param.Param) handler.Handler {
	return func(c *context.Context) *response.Response {
		args, err := handler.ExtractArgs(c.Args(), id)
		if err != nil {
			return response.BadRequest("invalid arguments: %v", err)
		}
		v, err := c.Store.GetCompositeSingleton(singleton, args[0])
		if err != nil {
			return response.Error("getting %s singleton %s: %v", s.Name(), singleton.Name, err)
		}
		if v == nil {
			return response.NotFoundWithMessage("%s identified with %v not found", s.Name(), args[0])
		}
		return response.OK(v)
	}
}

func GetCollectionHandler(s *store.Schema, collection *store.Collection, id param.Param) handler.Handler {
	return func(c *context.Context) *response.Response {
		args, err := handler.ExtractArgs(c.Args(), id)
		if err != nil {
			return response.BadRequest("invalid arguments: %v", err)
		}
		v, err := c.Store.GetCompositeCollection(collection, args[0])
		if err != nil {
			return response.Error("getting %s collection %s: %v", s.Name(), collection.Name, err)
		}
		if v == nil {
			return response.NotFoundWithMessage("%s identified with %v not found or its %s collection was empty", s.Name(), args[0], collection.Name)
		}
		return response.OK(v)
	}
}

func GetCollectionItemHandler(s *store.Schema, collection *store.Collection, id param.Param) handler.Handler {
	return func(c *context.Context) *response.Response {
		args, err := handler.ExtractArgs(c.Args(), id, param.String)
		if err != nil {
			return response.BadRequest("invalid arguments: %v", err)
		}
		id := args[0]
		item := args[1].(string)
		v, err := c.Store.GetCompositeCollectionItem(collection, id, item)
		if err != nil {
			return response.Error("getting %s collection %s item %s: %v", s.Name(), collection.Name, item, err)
		}
		if v == nil {
			return response.NotFoundWithMessage("%s identified with %v not found or its %s collection item %q was not found", s.Name(), id, collection.Name, item)
		}
		return response.OK(v)
	}
}
