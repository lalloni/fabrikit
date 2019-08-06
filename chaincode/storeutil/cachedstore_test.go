package storeutil

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lalloni/fabrikit/chaincode/store"
	"github.com/lalloni/fabrikit/chaincode/store/key"
	"github.com/lalloni/fabrikit/chaincode/test"
)

type Thing struct {
	ID   uint64
	Name string
}

var ThingSchema = store.MustPrepare(store.Composite{
	Name:            "Thing",
	Creator:         func() interface{} { return &Thing{} },
	IdentifierField: "ID",
	KeyBaseName:     "t",
	KeepRoot:        true,
	IdentifierKey: func(id interface{}) (*key.Key, error) {
		return key.NewBase("t", strconv.FormatUint(id.(uint64), 10)), nil
	},
})

type spyStore struct {
	invocations int
	store       store.Store
}

func (ss *spyStore) PutComposite(s *store.Schema, val interface{}) error {
	ss.invocations++
	return ss.store.PutComposite(s, val)
}

func (ss *spyStore) GetComposite(s *store.Schema, id interface{}) (interface{}, error) {
	ss.invocations++
	return ss.store.GetComposite(s, id)
}

func (ss *spyStore) HasComposite(s *store.Schema, id interface{}) (bool, error) {
	ss.invocations++
	return ss.store.HasComposite(s, id)
}

func (ss *spyStore) DelComposite(s *store.Schema, id interface{}) error {
	ss.invocations++
	return ss.store.DelComposite(s, id)
}

func (ss *spyStore) GetCompositeAll(s *store.Schema) ([]interface{}, error) {
	ss.invocations++
	return ss.store.GetCompositeAll(s)
}

func (ss *spyStore) GetCompositeRange(s *store.Schema, r *store.Range) ([]interface{}, error) {
	ss.invocations++
	return ss.store.GetCompositeRange(s, r)
}

func (ss *spyStore) DelCompositeRange(s *store.Schema, r *store.Range) ([]interface{}, error) {
	ss.invocations++
	return ss.store.DelCompositeRange(s, r)
}

func (ss *spyStore) PutCompositeSingleton(s *store.Singleton, id interface{}, val interface{}) error {
	ss.invocations++
	return ss.store.PutCompositeSingleton(s, id, val)
}

func (ss *spyStore) GetCompositeSingleton(s *store.Singleton, id interface{}) (interface{}, error) {
	ss.invocations++
	return ss.store.GetCompositeSingleton(s, id)
}

func (ss *spyStore) PutCompositeCollection(c *store.Collection, id interface{}, col interface{}) error {
	ss.invocations++
	return ss.store.PutCompositeCollection(c, id, col)
}

func (ss *spyStore) GetCompositeCollection(c *store.Collection, id interface{}) (interface{}, error) {
	ss.invocations++
	return ss.store.GetCompositeCollection(c, id)
}

func (ss *spyStore) PutCompositeCollectionItem(c *store.Collection, id interface{}, itemid string, val interface{}) error {
	ss.invocations++
	return ss.store.PutCompositeCollectionItem(c, id, itemid, val)
}

func (ss *spyStore) GetCompositeCollectionItem(c *store.Collection, id interface{}, itemid string) (interface{}, error) {
	ss.invocations++
	return ss.store.GetCompositeCollectionItem(c, id, itemid)
}

func (ss *spyStore) DelCompositeCollectionItem(c *store.Collection, id interface{}, itemid string) error {
	ss.invocations++
	return ss.store.DelCompositeCollectionItem(c, id, itemid)
}

func (ss *spyStore) PutValue(k *key.Key, val interface{}) error {
	ss.invocations++
	return ss.store.PutValue(k, val)
}

func (ss *spyStore) GetValue(k *key.Key, val interface{}) (bool, error) {
	ss.invocations++
	return ss.store.GetValue(k, val)
}

func (ss *spyStore) HasValue(k *key.Key) (bool, error) {
	ss.invocations++
	return ss.store.HasValue(k)
}

func (ss *spyStore) DelValue(k *key.Key) error {
	ss.invocations++
	return ss.store.DelValue(k)
}

func TestCachedStore(t *testing.T) {
	a := assert.New(t)

	mock := test.NewMock("test", nil)

	spy := &spyStore{store: store.New(mock)}
	cs := ReadCacheStore(spy)

	test.InTransaction(mock, func(tx string) {
		a.NoError(cs.PutComposite(ThingSchema, &Thing{ID: 1, Name: "foo"}))
		a.NoError(cs.PutComposite(ThingSchema, &Thing{ID: 2, Name: "bar"}))
		a.NoError(cs.PutComposite(ThingSchema, &Thing{ID: 3, Name: "baz"}))
	})

	invocations := spy.invocations
	thing, err := cs.GetComposite(ThingSchema, uint64(1))
	a.NoError(err)
	a.EqualValues(&Thing{ID: 1, Name: "foo"}, thing)
	a.EqualValues(invocations+1, spy.invocations)

	invocations = spy.invocations
	thing, err = cs.GetComposite(ThingSchema, uint64(1))
	a.NoError(err)
	a.EqualValues(&Thing{ID: 1, Name: "foo"}, thing)
	a.EqualValues(invocations, spy.invocations) // no new invocations

	invocations = spy.invocations
	thing, err = cs.GetComposite(ThingSchema, uint64(2))
	a.NoError(err)
	a.EqualValues(&Thing{ID: 2, Name: "bar"}, thing)
	a.EqualValues(invocations+1, spy.invocations)

	invocations = spy.invocations
	thing, err = cs.GetComposite(ThingSchema, uint64(2))
	a.NoError(err)
	a.EqualValues(&Thing{ID: 2, Name: "bar"}, thing)
	a.EqualValues(invocations, spy.invocations)

	invocations = spy.invocations
	thing, err = cs.GetComposite(ThingSchema, uint64(1))
	a.NoError(err)
	a.EqualValues(&Thing{ID: 1, Name: "foo"}, thing)
	a.EqualValues(invocations, spy.invocations)

}
