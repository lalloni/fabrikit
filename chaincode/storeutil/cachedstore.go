package storeutil

import (
	"github.com/mitchellh/hashstructure"

	"github.com/lalloni/fabrikit/chaincode/store"
	"github.com/lalloni/fabrikit/chaincode/store/key"
)

type op int

const (
	// cached operations
	getComposite op = iota
	getCompositeAll
	getCompositeRange
	getCompositeSingleton
	getCompositeCollection
	getCompositeCollectionItem
	hasComposite
	hasValue
)

func ReadCacheStore(s store.Store) store.Store {
	return &cachedstore{
		store: s,
		cache: newcache(),
	}
}

type cachedstore struct {
	store store.Store
	cache *cache
}

func (cs *cachedstore) PutComposite(s *store.Schema, val interface{}) error {
	return cs.store.PutComposite(s, val)
}

func (cs *cachedstore) GetComposite(s *store.Schema, id interface{}) (interface{}, error) {
	return cs.cache.Cached(request(getComposite, id), func() (interface{}, error) {
		return cs.store.GetComposite(s, id)
	})
}

func (cs *cachedstore) HasComposite(s *store.Schema, id interface{}) (bool, error) {
	v, err := cs.cache.Cached(request(hasComposite, id), func() (interface{}, error) {
		return cs.store.HasComposite(s, id)
	})
	return v.(bool), err
}

func (cs *cachedstore) DelComposite(s *store.Schema, id interface{}) error {
	return cs.store.DelComposite(s, id)
}

func (cs *cachedstore) GetCompositeAll(s *store.Schema) ([]interface{}, error) {
	v, err := cs.cache.Cached(request(getCompositeAll), func() (interface{}, error) {
		return cs.store.GetCompositeAll(s)
	})
	return v.([]interface{}), err
}

func (cs *cachedstore) GetCompositeRange(s *store.Schema, r *store.Range) ([]interface{}, error) {
	v, err := cs.cache.Cached(request(getCompositeRange, r), func() (interface{}, error) {
		return cs.store.GetCompositeRange(s, r)
	})
	return v.([]interface{}), err
}

func (cs *cachedstore) DelCompositeRange(s *store.Schema, r *store.Range) ([]interface{}, error) {
	return cs.store.DelCompositeRange(s, r)
}

func (cs *cachedstore) PutCompositeSingleton(s *store.Singleton, id interface{}, val interface{}) error {
	return cs.store.PutCompositeSingleton(s, id, val)
}

func (cs *cachedstore) GetCompositeSingleton(s *store.Singleton, id interface{}) (interface{}, error) {
	return cs.cache.Cached(request(getCompositeSingleton, id), func() (interface{}, error) {
		return cs.store.GetCompositeSingleton(s, id)
	})
}

func (cs *cachedstore) PutCompositeCollection(c *store.Collection, id interface{}, col interface{}) error {
	return cs.store.PutCompositeCollection(c, id, col)
}

func (cs *cachedstore) GetCompositeCollection(c *store.Collection, id interface{}) (interface{}, error) {
	return cs.cache.Cached(request(getCompositeCollection, id), func() (interface{}, error) {
		return cs.store.GetCompositeCollection(c, id)
	})
}

func (cs *cachedstore) PutCompositeCollectionItem(c *store.Collection, id interface{}, itemid string, val interface{}) error {
	return cs.store.PutCompositeCollectionItem(c, id, itemid, val)
}

func (cs *cachedstore) GetCompositeCollectionItem(c *store.Collection, id interface{}, itemid string) (interface{}, error) {
	return cs.cache.Cached(request(getCompositeCollectionItem, struct {
		i interface{}
		s string
	}{id, itemid}), func() (interface{}, error) {
		return cs.store.GetCompositeCollectionItem(c, id, itemid)
	})
}

func (cs *cachedstore) DelCompositeCollectionItem(c *store.Collection, id interface{}, itemid string) error {
	return cs.store.DelCompositeCollectionItem(c, id, itemid)
}

func (cs *cachedstore) PutValue(k *key.Key, val interface{}) error {
	return cs.store.PutValue(k, val)
}

func (cs *cachedstore) GetValue(k *key.Key, val interface{}) (bool, error) {
	// can't cache because output is passed-in as argument (val)
	return cs.store.GetValue(k, val)
}

func (cs *cachedstore) HasValue(k *key.Key) (bool, error) {
	v, err := cs.cache.Cached(request(hasValue, k.String()), func() (interface{}, error) {
		return cs.store.HasValue(k)
	})
	return v.(bool), err
}

func (cs *cachedstore) DelValue(k *key.Key) error {
	return cs.store.DelValue(k)
}

// request generates fast hashable structures for 0 to 3 values
// and uses (more expensive) hashstructure library for other cases
func request(o op, values ...interface{}) interface{} {
	switch len(values) {
	case 0:
		return o
	case 1:
		return struct {
			o op
			v interface{}
		}{o, values[0]}
	case 2:
		return struct {
			o  op
			v1 interface{}
			v2 interface{}
		}{o, values[0], values[1]}
	case 3:
		return struct {
			o  op
			v1 interface{}
			v2 interface{}
			v3 interface{}
		}{o, values[0], values[1], values[2]}
	default: // pay the price
		hash, err := hashstructure.Hash(struct {
			Op     op
			Values []interface{}
		}{o, values}, nil)
		if err != nil {
			// should never happen
			panic(err)
		}
		return hash
	}
}
