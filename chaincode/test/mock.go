package test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"

	"github.com/lalloni/fabrikit/chaincode"
	"github.com/lalloni/fabrikit/chaincode/context"
	"github.com/lalloni/fabrikit/chaincode/response"
	"github.com/lalloni/fabrikit/chaincode/router"
)

func NewMock(name string, r router.Router) *shim.MockStub {
	return shim.NewMockStub(name, chaincode.New(name, "test", r))
}

func MockTransactionStart(t *testing.T, stub *shim.MockStub) string {
	tx := uuid.New().String()
	stub.MockTransactionStart(tx)
	return tx
}

func MockTransactionEnd(t *testing.T, stub *shim.MockStub, tx string) {
	stub.MockTransactionEnd(tx)
}

func MockInvoke(t *testing.T, stub *shim.MockStub, function string, args ...interface{}) (string, *peer.Response, *response.Payload, error) {
	tx := uuid.New().String()
	aa := append([]interface{}{function}, args...)
	bs, err := arguments(aa)
	if err != nil {
		return "", nil, nil, err
	}
	f, _, err := context.ParseFunction([]byte(function))
	if err != nil {
		return "", nil, nil, err
	}
	res, payload, err := result(stub.MockInvoke(tx, bs))
	t.Logf("\n→ call function %q arguments: %s\n← response status: %v message: %q payload: %v error: %v", f, format(skip(1, bs)), res.Status, res.Message, strings.Trim(string(res.Payload), "\n"), err)
	return tx, res, payload, err
}

func MockInit(t *testing.T, stub *shim.MockStub, args ...interface{}) (string, *peer.Response, *response.Payload, error) {
	tx := uuid.New().String()
	bs, err := arguments(args)
	if err != nil {
		return "", nil, nil, err
	}
	res, payload, err := result(stub.MockInit(tx, bs))
	t.Logf("\n→ call init arguments: %s\n← response status: %v message: %q payload: %v error: %v", format(skip(1, bs)), res.Status, res.Message, strings.Trim(string(res.Payload), "\n"), err)
	return tx, res, payload, err
}

func result(r peer.Response) (*peer.Response, *response.Payload, error) {
	p := response.Payload{}
	if len(r.Payload) == 0 {
		return &r, &p, nil
	}
	err := json.Unmarshal(r.Payload, &p)
	if err != nil {
		return nil, nil, err
	}
	return &r, &p, nil
}

func arguments(args []interface{}) ([][]byte, error) {
	bs := [][]byte{}
	for _, arg := range args {
		switch v := arg.(type) {
		case []byte:
			bs = append(bs, v)
		case string:
			bs = append(bs, []byte(v))
		default:
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			bs = append(bs, b)
		}
	}
	return bs, nil
}

func skip(n int, bss [][]byte) [][]byte {
	if len(bss) >= n {
		return bss[n:]
	}
	return bss
}

func format(bss [][]byte) string {
	ss := []string{}
	for _, bs := range bss {
		ss = append(ss, fmt.Sprintf("%#v", string(bs)))
	}
	return strings.Join(ss, ",")
}
