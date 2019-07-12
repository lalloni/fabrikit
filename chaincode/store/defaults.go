package store

import (
	"github.com/lalloni/fabrikit/chaincode/store/filtering"
	"github.com/lalloni/fabrikit/chaincode/store/marshaling"
)

var (

	// DefaultMarshaling es el marshaling por defecto
	DefaultMarshaling = marshaling.JSON()

	// DefaultFiltering es el filtering por defecto
	DefaultFiltering = filtering.Copy()
)
