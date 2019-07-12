package chaincode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lalloni/fabrikit/chaincode/authorization"
	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/response"
	"github.com/lalloni/fabrikit/chaincode/response/status"
	"github.com/lalloni/fabrikit/chaincode/router"
	"github.com/lalloni/fabrikit/chaincode/test"
)

func TestInitHandler(t *testing.T) {
	a := assert.New(t)
	init := false
	r := router.New()
	r.SetInitHandler(nil, func(*context.Context) *response.Response {
		init = true
		return response.OK("blah!")
	})
	mock := test.NewMock("cc", r)
	_, res, p, err := test.MockInit(t, mock)
	a.NoError(err)
	a.True(init)
	a.EqualValues(status.OK, res.Status)
	a.EqualValues("blah!", p.Content)
}

func TestInvokeAuthorizationHandler(t *testing.T) {
	a := assert.New(t)

	var call bool
	r := router.New()
	r.SetHandler("f1", authorization.Allowed, func(*context.Context) *response.Response {
		call = true
		return response.OK("blah!")
	})
	r.SetHandler("f2", authorization.Forbidden, func(*context.Context) *response.Response {
		call = true
		return response.OK("bleh!") // never invoked
	})
	mock := test.NewMock("cc", r)

	call = false
	_, res, p, err := test.MockInvoke(t, mock, "f1")
	a.NoError(err)
	a.True(call)
	a.EqualValues(status.OK, res.Status)
	a.EqualValues("blah!", p.Content)

	call = false
	_, res, _, err = test.MockInvoke(t, mock, "f2")
	a.NoError(err)
	a.False(call)
	a.EqualValues(status.Forbidden, res.Status)
}

func TestDebugCall(t *testing.T) {
	a := assert.New(t)

	r := router.New()
	r.SetHandler("success", nil, func(*context.Context) *response.Response {
		return response.OK("blah!")
	})
	r.SetHandler("fail", nil, func(*context.Context) *response.Response {
		return response.Error("bleh!")
	})
	mock := test.NewMock("cc", r)

	_, res, p, err := test.MockInvoke(t, mock, "success")
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	a.EqualValues(&response.Payload{Content: "blah!"}, p)
	a.EqualValues("", res.Message)

	_, res, _, err = test.MockInvoke(t, mock, "fail")
	a.NoError(err)
	a.Nil(res.Payload)
	a.EqualValues(status.Error, res.Status)
	a.EqualValues("bleh!", res.Message)

	tx, res, p, err := test.MockInvoke(t, mock, "success?debug")
	a.NoError(err)
	a.EqualValues(status.OK, res.Status)
	a.EqualValues(&response.Payload{
		Content:     "blah!",
		Chaincode:   &response.Chaincode{Version: "test"},
		Transaction: &response.Transaction{ID: tx, Function: "success"},
	}, p)
	a.EqualValues("", res.Message)

	tx, res, p, err = test.MockInvoke(t, mock, "fail?debug")
	a.NoError(err)
	a.EqualValues(status.Error, res.Status)
	a.EqualValues(&response.Payload{
		Chaincode:   &response.Chaincode{Version: "test"},
		Transaction: &response.Transaction{ID: tx, Function: "fail"},
	}, p)
	a.EqualValues("bleh!", res.Message)

}
