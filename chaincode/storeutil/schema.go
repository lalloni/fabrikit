package storeutil

import (
	"strconv"

	"github.com/lalloni/fabrikit/chaincode/store"
	"github.com/lalloni/fabrikit/chaincode/store/key"
)

func Uint64Key(name string) store.KeyFunc {
	return func(id interface{}) (*key.Key, error) {
		return key.NewBase(name, strconv.FormatUint(id.(uint64), 10)), nil
	}
}

func Uint64Identifier(seg int) store.ValFunc {
	return func(k *key.Key) (interface{}, error) {
		return strconv.ParseUint(k.Base[seg].Value, 10, 64)
	}
}
