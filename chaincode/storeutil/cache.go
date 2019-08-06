package storeutil

import (
	"github.com/lalloni/fabrikit/chaincode/logging"
)

const CacheSizeWarningThreshold = 1000

var log = logging.ChaincodeLogger("fabrikit-cache")

func newcache() *cache {
	return &cache{
		data: make(map[interface{}]response),
	}
}

type cache struct {
	data map[interface{}]response
}

func (c *cache) Cached(key interface{}, source func() (interface{}, error)) (interface{}, error) {
	r, ok := c.data[key]
	if ok {
		return r.value, r.error
	}
	r.value, r.error = source()
	if r.error == nil {
		c.data[key] = r
		l := len(c.data)
		if l > CacheSizeWarningThreshold {
			log.Warningf("cache size of %v exceeded warning threshold of %v", l, CacheSizeWarningThreshold)
		}
	}
	return r.value, r.error
}

type response struct {
	value interface{}
	error error
}
