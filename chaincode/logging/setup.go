package logging

import (
	"os"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func Setup(name string) *shim.ChaincodeLogger {
	log := shim.NewLogger(name)
	l := os.Getenv("SHIM_LOGGING_LEVEL")
	if l != "" {
		level, err := shim.LogLevel(os.Getenv("SHIM_LOGGING_LEVEL"))
		if err != nil {
			log.Errorf("parsing SHIM_LOGGING_LEVEL from value %q: %v", l, err)
			os.Exit(1)
		}
		shim.SetLoggingLevel(level)
	}
	return log
}
