package authorization

import (
	"github.com/lalloni/fabrikit/chaincode/context"
)

type Check func(*context.Context) error
