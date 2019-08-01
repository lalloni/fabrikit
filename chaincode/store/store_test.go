package store_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"

	"github.com/lalloni/fabrikit/chaincode/store"
	"github.com/lalloni/fabrikit/chaincode/store/key"
)

type Thingy struct {
	Blah string `json:"blah,omitempty"`
}

type Thing struct {
	ID       uint64   `json:"id,omitempty"`
	Name     string   `json:"name,omitempty"`
	Age      uint8    `json:"age,omitempty"`
	Thingies []Thingy `json:"thingies,omitempty"`
	Error    string   `json:"error,omitempty"`
}

type Other struct {
	Name   string `json:"name,omitempty"`
	Number int    `json:"number,omitempty"`
}

type Foo struct {
	Some  string `json:"some,omitempty"`
	Num   int    `json:"num,omitempty"`
	Error string `json:"error,omitempty"`
}

type Item struct {
	Name     string  `json:"name,omitempty"`
	Quantity float64 `json:"quantity,omitempty"`
}

type Compo struct {
	Thing  *Thing           `json:"thing,omitempty"`
	Other  *Other           `json:"other,omitempty"`
	Items  map[string]*Item `json:"items,omitempty"`
	Foos   map[string]*Foo  `json:"foos,omitempty"`
	Errors interface{}      `json:"errors,omitempty"`
	Name   string           `json:"name,omitempty"`
}

var cc = store.MustPrepare(store.Composite{
	Name:        "compo",
	KeepRoot:    true,
	Creator:     func() interface{} { return &Compo{Items: map[string]*Item{}, Foos: map[string]*Foo{}} },
	KeyBaseName: "compo",
	IdentifierGetter: func(v interface{}) interface{} {
		return v.(*Compo).Thing.ID
	},
	IdentifierKey: func(id interface{}) (*key.Key, error) {
		return key.NewBase("compo", strconv.FormatUint(id.(uint64), 10)), nil
	},
	KeyIdentifier: func(k *key.Key) (interface{}, error) {
		return strconv.ParseUint(k.Base[0].Value, 10, 64)
	},
	Singletons: []store.Singleton{
		{Tag: "thing",
			Field:   "Thing",
			Creator: func() interface{} { return &Thing{} },
		},
		{Tag: "other",
			Field: "Other",
		},
	},
	Collections: []store.Collection{
		{Tag: "item",
			Field:     "Items",
			Creator:   func() interface{} { return &Item{} },
			Collector: func(v interface{}, i store.Item) { v.(*Compo).Items[i.Identifier] = i.Value.(*Item) },
			Enumerator: func(v interface{}) []store.Item {
				items := []store.Item{}
				for k, v := range v.(*Compo).Items {
					items = append(items, store.Item{Identifier: k, Value: v})
				}
				return items
			},
		},
		{Tag: "foos", Field: "Foos"},
	},
})

func TestPutAndGetValue(t *testing.T) {
	shim.SetLoggingLevel(shim.LogDebug)
	logging.SetLevel(logging.DEBUG, "mock")
	a := assert.New(t)

	stub := shim.NewMockStub("test", nil)

	st := store.New(stub)
	key := key.NewBase("thing", "1")

	stub.MockTransactionStart("x")

	t1 := &Thing{1234, "PP", 16, []Thingy{{"A"}, {"B"}}, ""}
	err := st.PutValue(key, t1)
	a.NoError(err)

	stub.MockTransactionEnd("x")

	t2 := &Thing{1234, "AA", 100, []Thingy{}, ""}
	_, err = st.GetValue(key, t2)
	a.NoError(err)
	a.Equal(t1, t2)
}

func TestPutAndGetComposite(t *testing.T) {
	a := assert.New(t)

	shim.SetLoggingLevel(shim.LogDebug)
	logging.SetLevel(logging.DEBUG, "mock")

	stub := shim.NewMockStub("test", nil)
	st := store.New(stub)

	c1 := &Compo{
		Thing: &Thing{1234, "PP", 16, []Thingy{{"A"}, {"B"}}, ""},
		Other: &Other{"TT", 2123},
		Items: map[string]*Item{
			"a": {Name: "Pedro", Quantity: 10.0},
			"b": {Name: "Pablo", Quantity: 20.0},
		},
		Foos: map[string]*Foo{
			"foo1": {Some: "bar", Num: 634},
			"foo2": {Some: "baz", Num: 634},
		},
		Name: "some very interesting name",
	}

	stub.MockTransactionStart("x")
	err := st.PutComposite(cc, c1)
	stub.MockTransactionEnd("x")
	a.NoError(err)
	t.Logf("put: %+v", mustMarshal(c1))

	c2, err := st.GetComposite(cc, c1.Thing.ID)
	a.NoError(err)
	t.Logf("get: %+v", mustMarshal(c2))
	a.Equal(c1, c2)
}

func TestGetMissingComposite(t *testing.T) {
	a := assert.New(t)

	shim.SetLoggingLevel(shim.LogDebug)
	logging.SetLevel(logging.DEBUG, "mock")

	stub := shim.NewMockStub("test", nil)
	st := store.New(stub)

	c, err := st.GetComposite(cc, uint64(1))
	a.NoError(err)
	a.Nil(c)
}

func TestPutAndDeleteComposite(t *testing.T) {
	a := assert.New(t)

	shim.SetLoggingLevel(shim.LogDebug)
	logging.SetLevel(logging.DEBUG, "mock")

	stub := shim.NewMockStub("test", nil)
	st := store.New(stub)

	c1 := &Compo{
		Thing: &Thing{1234, "PP", 16, []Thingy{{"A"}, {"B"}}, ""},
		Items: map[string]*Item{
			"a": {
				Name:     "Pedro",
				Quantity: 10.0,
			},
			"b": {
				Name:     "Pablo",
				Quantity: 20.0,
			},
		},
	}

	stub.MockTransactionStart("x")
	err := st.PutComposite(cc, c1)
	stub.MockTransactionEnd("x")
	a.NoError(err)
	t.Logf("put: %+v", mustMarshal(c1))

	has, err := st.HasComposite(cc, c1.Thing.ID)
	a.NoError(err)
	a.True(has)

	err = st.DelComposite(cc, c1.Thing.ID)
	a.NoError(err)

	has, err = st.HasComposite(cc, c1.Thing.ID)
	a.NoError(err)
	a.False(has)

}

func TestPutPartialComposite(t *testing.T) {
	a := assert.New(t)

	shim.SetLoggingLevel(shim.LogDebug)
	logging.SetLevel(logging.DEBUG, "mock")

	stub := shim.NewMockStub("test", nil)
	st := store.New(stub)

	id := uint64(1234)

	c1 := &Compo{
		Thing: &Thing{id, "PP", 16, []Thingy{{"A"}, {"B"}}, ""},
		Other: &Other{"TT", 2123},
		Items: map[string]*Item{
			"a": {Name: "Pedro", Quantity: 10.0},
			"b": {Name: "Pablo", Quantity: 20.0},
		},
		Foos: map[string]*Foo{
			"foo1": {Some: "bar", Num: 634},
			"foo2": {Some: "baz", Num: 634},
		},
	}

	stub.MockTransactionStart("x")
	err := st.PutComposite(cc, c1)
	stub.MockTransactionEnd("x")
	a.NoError(err)
	t.Logf("put: %s", mustMarshal(c1))

	c2, err := st.GetComposite(cc, id)
	a.NoError(err)
	t.Logf("get: %s", mustMarshal(c2))
	a.Equal(c1, c2)

	c3 := *c1
	c3.Other = nil
	stub.MockTransactionStart("x")
	err = st.PutComposite(cc, &c3)
	stub.MockTransactionEnd("x")
	a.NoError(err)
	t.Logf("put: %s", mustMarshal(&c3))

	c2, err = st.GetComposite(cc, id)
	a.NoError(err)
	t.Logf("get: %s", mustMarshal(c2))
	a.Equal(c1, c2)
	a.NotNil(c2.(*Compo).Other)

}

func TestMemberError(t *testing.T) {

	a := assert.New(t)

	shim.SetLoggingLevel(shim.LogDebug)
	logging.SetLevel(logging.DEBUG, "mock")

	stub := shim.NewMockStub("test", nil)

	id := uint64(1)

	stub.MockTransactionStart("x")

	err := stub.PutState("compo:1#thing", []byte(`{"id":1,"age":"un montón","name":"pedro"}`))
	a.NoError(err)
	err = stub.PutState("compo:1#foos:1", []byte(`{"some":"bla bla","num":"cuarenta y dos"}`))
	a.NoError(err)
	err = stub.PutState("compo:1#item:a", []byte(`{"name":"algo","quantity":1.23}`))
	a.NoError(err)
	err = stub.PutState("compo:1#wit", []byte(`1`))
	a.NoError(err)
	stub.MockTransactionEnd("x")

	st1 := store.New(stub)
	_, err = st1.GetComposite(cc, id)
	a.NoError(err)

	st2 := store.New(stub, store.SetErrors(true))
	c2, err := st2.GetComposite(cc, id)
	a.NoError(err)
	t.Logf("got: %s", mustMarshal(c2))
	c := c2.(*Compo)
	a.NotEmpty(c.Thing.Error)
	a.NotEmpty(c.Foos)
	a.NotEmpty(c.Foos["1"].Error)
	a.NotEmpty(c.Errors)

}

func TestDelCompositeRange(t *testing.T) {
	a := assert.New(t)

	shim.SetLoggingLevel(shim.LogDebug)
	logging.SetLevel(logging.DEBUG, "mock")

	stub := shim.NewMockStub("test", nil)
	st := store.New(stub)

	c1 := &Compo{
		Thing: &Thing{1234, "PP", 16, []Thingy{{"A"}, {"B"}}, ""},
		Items: map[string]*Item{"a": {Name: "Pedro", Quantity: 10.0}},
	}

	for id := 100; id < 110; id++ {
		c1.Thing.ID = uint64(id)
		stub.MockTransactionStart("x-" + strconv.Itoa(id))
		err := st.PutComposite(cc, c1)
		stub.MockTransactionEnd("x-" + strconv.Itoa(id))
		a.NoError(err)
		t.Logf("put: %s", mustMarshal(c1))
	}

	ids, err := st.DelCompositeRange(cc, &store.Range{First: uint64(102), Last: uint64(105)})
	a.NoError(err)
	a.Len(ids, 4)
	t.Logf("deleted: %s", mustMarshal(ids))

	has, err := st.HasComposite(cc, uint64(101))
	a.NoError(err)
	a.True(has)

	has, err = st.HasComposite(cc, uint64(102))
	a.NoError(err)
	a.False(has)

	has, err = st.HasComposite(cc, uint64(105))
	a.NoError(err)
	a.False(has)

	has, err = st.HasComposite(cc, uint64(106))
	a.NoError(err)
	a.True(has)

}

func TestGetCompositeRange(t *testing.T) {
	a := assert.New(t)

	shim.SetLoggingLevel(shim.LogDebug)
	logging.SetLevel(logging.DEBUG, "mock")

	stub := shim.NewMockStub("test", nil)
	st := store.New(stub)

	c1 := &Compo{
		Thing: &Thing{1234, "PP", 16, []Thingy{{"A"}, {"B"}}, ""},
		Items: map[string]*Item{"a": {Name: "Pedro", Quantity: 10.0}},
		Foos:  map[string]*Foo{},
	}

	i0, i1 := 102, 105
	k := []*Compo{}
	for id := 100; id < 110; id++ {
		c1.Thing.ID = uint64(id)
		stub.MockTransactionStart("x-" + strconv.Itoa(id))
		err := st.PutComposite(cc, c1)
		stub.MockTransactionEnd("x-" + strconv.Itoa(id))
		a.NoError(err)
		t.Logf("put: %s", mustMarshal(c1))
		if i0 <= id && id <= i1 {
			o := &Compo{}
			a.NoError(deepCopy(c1, o))
			k = append(k, o)
		}
	}
	t.Logf("kept: %s", mustMarshal(k))

	cs, err := st.GetCompositeRange(cc, &store.Range{First: uint64(i0), Last: uint64(i1)})
	a.NoError(err)
	a.Len(cs, i1-i0+1)
	t.Logf("got: %s", mustMarshal(cs))

	for i := 0; i <= i1-i0; i++ {
		a.EqualValues(k[i], cs[i])
	}

}

func TestGetCompositeAll(t *testing.T) {
	a := assert.New(t)

	shim.SetLoggingLevel(shim.LogDebug)
	logging.SetLevel(logging.DEBUG, "mock")

	stub := shim.NewMockStub("test", nil)
	st := store.New(stub)

	c1 := &Compo{
		Thing: &Thing{1234, "PP", 16, []Thingy{{"A"}, {"B"}}, ""},
		Items: map[string]*Item{"a": {Name: "Pedro", Quantity: 10.0}},
		Foos:  map[string]*Foo{},
	}

	for id := 100; id < 110; id++ {
		c1.Thing.ID = uint64(id)
		stub.MockTransactionStart("x-" + strconv.Itoa(id))
		err := st.PutComposite(cc, c1)
		stub.MockTransactionEnd("x-" + strconv.Itoa(id))
		a.NoError(err)
	}

	rs, err := st.GetCompositeAll(cc)
	a.NoError(err)
	idx := map[uint64]*Compo{}
	for _, r := range rs {
		c := r.(*Compo)
		idx[c.Thing.ID] = c
	}
	a.Len(idx, 10)
	for id := 100; id < 110; id++ {
		a.NotNil(idx[uint64(id)])
	}
}

func TestGetCompositeSingleton(t *testing.T) {

	a := assert.New(t)

	c1 := &Compo{
		Thing: &Thing{1234, "PP", 16, []Thingy{{"A"}, {"B"}}, ""},
	}

	shim.SetLoggingLevel(shim.LogDebug)
	logging.SetLevel(logging.DEBUG, "mock")

	stub := shim.NewMockStub("test", nil)
	st := store.New(stub)

	stub.MockTransactionStart("x")
	err := st.PutComposite(cc, c1)
	stub.MockTransactionEnd("x")
	a.NoError(err)

	thing, err := st.GetCompositeSingleton(cc.Singleton("thing"), c1.Thing.ID)
	a.NoError(err)
	a.EqualValues(c1.Thing, thing)

}

func TestPutCompositeSingleton(t *testing.T) {

	a := assert.New(t)

	c1 := &Compo{
		Thing: &Thing{1234, "PP", 16, []Thingy{{"A"}, {"B"}}, ""},
	}

	shim.SetLoggingLevel(shim.LogDebug)
	logging.SetLevel(logging.DEBUG, "mock")

	stub := shim.NewMockStub("test", nil)
	st := store.New(stub)

	stub.MockTransactionStart("x")
	err := st.PutComposite(cc, c1)
	stub.MockTransactionEnd("x")
	a.NoError(err)

	thing2 := &Thing{1234, "QQQ", 32, []Thingy{{"X"}, {"Y"}, {"Z"}}, "blabla"}
	stub.MockTransactionStart("x")
	err = st.PutCompositeSingleton(cc.Singleton("thing"), c1.Thing.ID, thing2)
	stub.MockTransactionEnd("x")
	a.NoError(err)

	c2, err := st.GetComposite(cc, c1.Thing.ID)
	a.NoError(err)
	a.EqualValues(thing2, c2.(*Compo).Thing)

}

func mustMarshal(v interface{}) string {
	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

func deepCopy(src, tgt interface{}) error {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(src)
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewReader(b.Bytes())).Decode(tgt)
}
