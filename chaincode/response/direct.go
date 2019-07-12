package response

import (
	"github.com/hyperledger/fabric/protos/peer"
)

func Direct(res peer.Response) *Response {
	return StatusWithResult(-1, res)
}
