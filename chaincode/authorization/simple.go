package authorization

import (
	"github.com/pkg/errors"

	"github.com/lalloni/fabrikit/chaincode/context"
)

func Allowed(*context.Context) error {
	return nil
}

func Forbidden(*context.Context) error {
	return errors.New("not allowed")
}
