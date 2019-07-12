package logging

import (
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func ChaincodeLogger(path ...string) *shim.ChaincodeLogger {
	return shim.NewLogger(strings.Join(append([]string{"cc"}, path...), "/"))
}
