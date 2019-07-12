package router_test

import (
	"testing"

	"github.com/hyperledger/fabric/protos/peer"
	"github.com/stretchr/testify/assert"

	"github.com/lalloni/fabrikit/chaincode/authorization"
	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/response"
	"github.com/lalloni/fabrikit/chaincode/response/status"
	"github.com/lalloni/fabrikit/chaincode/router"
	"github.com/lalloni/fabrikit/chaincode/test"
)

func TestRouter(t *testing.T) {
	var (
		tx   string
		res  *peer.Response
		p    *response.Payload
		err  error
		pong bool
		ping = func(_ *context.Context) *response.Response {
			pong = true
			return response.OK("ok!")
		}
	)

	a := assert.New(t)
	r := router.New()
	stub := test.NewMock("test", r)

	_, res, p, err = test.MockInit(t, stub, nil)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	a.EqualValues(nil, p.Content)

	pong = false
	r.SetInitHandler(authorization.Allowed, ping)
	_, res, p, err = test.MockInit(t, stub, nil)
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	a.EqualValues("ok!", p.Content)
	a.True(pong)

	pong = false
	r.SetInitHandler(authorization.Forbidden, ping)
	_, res, p, err = test.MockInit(t, stub, nil)
	a.NoError(err)
	a.EqualValues(status.Forbidden, res.Status)
	a.False(pong)
	a.EqualValues(nil, p.Content)

	tx, res, p, err = test.MockInvoke(t, stub, "h")
	a.NoError(err)
	a.EqualValues(status.BadRequest, res.Status)
	a.NotEmpty(tx)
	a.Nil(p.Content)

	r.SetHandler("h", nil, nil)
	tx, res, p, err = test.MockInvoke(t, stub, "h")
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	a.NotEmpty(tx)
	a.Nil(p.Content)

	r.SetHandler("h", authorization.Forbidden, nil)
	tx, res, p, err = test.MockInvoke(t, stub, "h")
	a.NoError(err)
	a.EqualValues(status.Forbidden, res.Status)
	a.NotEmpty(tx)
	a.Nil(p.Content)

	pong = false
	r.SetHandler("h", authorization.Forbidden, ping)
	tx, res, p, err = test.MockInvoke(t, stub, "h")
	a.NoError(err)
	a.EqualValues(status.Forbidden, res.Status)
	a.False(pong)
	a.NotEmpty(tx)
	a.Nil(p.Content)

	pong = false
	r.SetHandler("h", nil, ping)
	tx, res, p, err = test.MockInvoke(t, stub, "h")
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	a.True(pong)
	a.NotEmpty(tx)
	a.EqualValues("ok!", p.Content)

	pong = false
	r.SetHandler("h", authorization.Allowed, ping)
	tx, res, p, err = test.MockInvoke(t, stub, "h")
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	a.True(pong)
	a.NotEmpty(tx)
	a.EqualValues("ok!", p.Content)
}
