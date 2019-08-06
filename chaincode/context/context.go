package context

import (
	"crypto/x509"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/pkg/errors"

	"github.com/lalloni/fabrikit/chaincode/logging"
	"github.com/lalloni/fabrikit/chaincode/store"
	"github.com/lalloni/fabrikit/chaincode/storeutil"
)

func New(stub shim.ChaincodeStubInterface, name, version string, path ...string) *Context {
	c := &Context{
		name:    name,
		version: version,
		path:    append([]string{name, version}, path...),
		Stub:    stub,
		Store:   store.New(stub),
	}
	args := stub.GetArgs()
	if len(args) > 0 {
		fun, opts, err := ParseFunction(args[0])
		if err != nil {
			c.Logger().Warningf("parsing function call options: %v", err)
		}
		c.function = fun
		c.options = opts
	}
	if c.options["storecache"] != "no" {
		// setup cached store unless invoked with storecache=no option
		c.Store = storeutil.ReadCacheStore(c.Store)
	}
	return c
}

type Context struct {
	Stub        shim.ChaincodeStubInterface
	Store       store.Store
	name        string
	version     string
	function    string
	options     map[string]string
	path        []string
	clientid    cid.ClientIdentity
	clientcrt   *x509.Certificate
	clientmspid string
}

func (ctx *Context) Version() string {
	return ctx.version
}

func (ctx *Context) Logger(path ...string) *shim.ChaincodeLogger {
	return logging.ChaincodeLogger(append(ctx.path, path...)...)
}

func (ctx *Context) ClientIdentity() (cid.ClientIdentity, error) {
	if ctx.clientid == nil {
		id, err := cid.New(ctx.Stub)
		if err != nil {
			return nil, errors.Wrap(err, "creating client identity")
		}
		ctx.clientid = id
	}
	return ctx.clientid, nil
}

func (ctx *Context) ClientCertificate() (*x509.Certificate, error) {
	if ctx.clientcrt == nil {
		id, err := ctx.ClientIdentity()
		if err != nil {
			return nil, err
		}
		cert, err := id.GetX509Certificate()
		if err != nil {
			return nil, errors.Wrap(err, "getting client certificate")
		}
		ctx.clientcrt = cert
	}
	return ctx.clientcrt, nil
}

func (ctx *Context) ClientMSPID() (string, error) {
	if ctx.clientmspid == "" {
		id, err := ctx.ClientIdentity()
		if err != nil {
			return "", err
		}
		mspid, err := id.GetMSPID()
		if err != nil {
			return "", errors.Wrap(err, "getting client mspid")
		}
		ctx.clientmspid = mspid
	}
	return ctx.clientmspid, nil
}

func (ctx *Context) Function() string {
	return ctx.function
}

func (ctx *Context) Option(name string) (string, bool) {
	v, present := ctx.options[name]
	return v, present
}

func (ctx *Context) Args() [][]byte {
	return ctx.Stub.GetArgs()[1:]
}

func (ctx *Context) ArgBytes(n int) ([]byte, error) {
	args := ctx.Stub.GetArgs()
	if len(args) < n+1 {
		return nil, errors.Errorf("argument %d is required", n)
	}
	return args[n], nil
}

func (ctx *Context) ArgString(n int) (string, error) {
	bs, err := ctx.ArgBytes(n)
	return string(bs), err
}

func (ctx *Context) ArgInt64(n int) (int64, error) {
	bs, err := ctx.ArgBytes(n)
	if err != nil {
		return 0, err
	}
	r, err := strconv.ParseInt(string(bs), 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "argument %d must be an integer", n)
	}
	return r, nil
}

func (ctx *Context) ArgUint64(n int) (uint64, error) {
	bs, err := ctx.ArgBytes(n)
	if err != nil {
		return 0, err
	}
	r, err := strconv.ParseUint(string(bs), 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "argument %d must be a natural integer", n)
	}
	return r, nil
}

func (ctx *Context) ArgKV(pos int) (string, []byte, error) {
	key, err := ctx.ArgString(pos)
	if err != nil {
		return "", nil, err
	}
	val, err := ctx.ArgBytes(pos + 1)
	if err != nil {
		return "", nil, err
	}
	return key, val, nil
}
