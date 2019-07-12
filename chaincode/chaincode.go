package chaincode

import (
	"bytes"
	"encoding/json"
	"unicode/utf8"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/handler"
	"github.com/lalloni/fabrikit/chaincode/logging"
	"github.com/lalloni/fabrikit/chaincode/response"
	"github.com/lalloni/fabrikit/chaincode/router"
)

func New(name string, version string, r router.Router) shim.Chaincode {
	log := logging.ChaincodeLogger(name)
	log.Info("created")
	res := &cc{
		name:    name,
		version: version,
		router:  r,
		log:     log,
	}
	return res
}

type cc struct {
	name    string
	version string
	router  router.Router
	log     *shim.ChaincodeLogger
}

func (c *cc) Init(stub shim.ChaincodeStubInterface) peer.Response {
	ctx := context.New(stub, c.name, c.version, "init")
	logger := ctx.Logger()
	logger.Debug("begin request processing")
	handle := c.router.InitHandler()
	if handle != nil {
		return c.response(ctx, logger, handle(ctx))
	}
	res := c.response(ctx, logger, response.OK(nil))
	logger.Debugf("end request processing with response status %q", res.GetStatus())
	return res
}

func (c *cc) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	ctx := context.New(stub, c.name, c.version, "invoke", stub.GetTxID())
	logger := ctx.Logger(ctx.Function())
	logger.Debug("begin request processing")
	handle := c.router.Handler(router.Name(ctx.Function()))
	if handle == nil {
		handle = handler.NotImplementedHandler
	}
	res := c.response(ctx, logger, handle(ctx))
	logger.Debugf("end request processing with response status %q", res.GetStatus())
	return res
}

func (c *cc) response(ctx *context.Context, logger *shim.ChaincodeLogger, r *response.Response) peer.Response {
	if r.Status < 0 {
		return r.Payload.Content.(peer.Response)
	}
	var payload []byte
	_, debug := ctx.Option("debug")
	if debug {
		if r.Payload == nil {
			r.Payload = &response.Payload{}
		}
		r.Payload.Chaincode = &response.Chaincode{Version: ctx.Version()}
		r.Payload.Transaction = &response.Transaction{ID: ctx.Stub.GetTxID(), Function: ctx.Function()}
		mspid, err := ctx.ClientMSPID()
		if err != nil {
			logger.Warningf("getting MSPID: %v", err)
		}
		var subject, issuer string
		cert, err := ctx.ClientCertificate()
		if err != nil {
			logger.Warningf("getting client certificate: %v", err)
		} else {
			subject = cert.Subject.String()
			issuer = cert.Issuer.String()
		}
		if mspid != "" || subject != "" || issuer != "" {
			r.Payload.Client = &response.Client{MSPID: mspid, Subject: subject, Issuer: issuer}
		}
	}
	if r.Payload != nil {
		if bs, ok := r.Payload.Content.([]byte); ok {
			if utf8.Valid(bs) {
				r.Payload.Content = string(bs)
			} else {
				// JSON encoding will encode []byte as a base64 string
				r.Payload.ContentEncoding = "base64"
			}
		}
		b := &bytes.Buffer{}
		enc := json.NewEncoder(b)
		enc.SetEscapeHTML(false) // do not html-escape "<", ">", "&"
		err := enc.Encode(r.Payload)
		if err != nil {
			return c.response(ctx, logger, response.Error("encoding response payload: %v", err))
		}
		payload = b.Bytes()
		// drop extra newline added by enc.Encode()
		p := len(payload) - 1
		if payload[p] == '\n' {
			payload = payload[0:p]
		}
	}
	return peer.Response{
		Status:  r.Status,
		Message: r.Message,
		Payload: payload,
	}
}
