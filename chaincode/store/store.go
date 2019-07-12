package store

import (
	"reflect"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"github.com/pkg/errors"

	"github.com/lalloni/fabrikit/chaincode/store/filtering"
	"github.com/lalloni/fabrikit/chaincode/store/key"
	"github.com/lalloni/fabrikit/chaincode/store/marshaling"
)

var log = shim.NewLogger("store")

type Store interface {
	PutValue(key *key.Key, val interface{}) error
	GetValue(key *key.Key, val interface{}) (bool, error)
	HasValue(key *key.Key) (bool, error)
	DelValue(key *key.Key) error

	PutComposite(s *Schema, val interface{}) error
	GetComposite(s *Schema, id interface{}) (interface{}, error)
	HasComposite(s *Schema, id interface{}) (bool, error)
	DelComposite(s *Schema, id interface{}) error
	GetCompositeAll(s *Schema) ([]interface{}, error)
	GetCompositeRange(s *Schema, r *Range) ([]interface{}, error)
	DelCompositeRange(s *Schema, r *Range) ([]interface{}, error)
}

func New(stub shim.ChaincodeStubInterface, opts ...Option) Store {
	s := &simplestore{
		stub:       stub,
		marshaling: DefaultMarshaling,
		filtering:  DefaultFiltering,
		sep:        key.DefaultSep,
		log:        shim.NewLogger("store"),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type simplestore struct {
	stub       shim.ChaincodeStubInterface
	log        *shim.ChaincodeLogger
	marshaling marshaling.Marshaling
	filtering  filtering.Filtering
	sep        *key.Sep
	seterrs    bool
}

func (ss *simplestore) PutValue(k *key.Key, value interface{}) error {
	return ss.internalPutValue(k, value)
}

func (ss *simplestore) GetValue(k *key.Key, value interface{}) (bool, error) {
	return ss.internalGetValue(k, value)
}

func (ss *simplestore) HasValue(k *key.Key) (bool, error) {
	return ss.internalHasValue(k)
}

func (ss *simplestore) DelValue(k *key.Key) error {
	return ss.internalDelValue(k)
}

func (ss *simplestore) PutComposite(s *Schema, val interface{}) error {
	we, err := s.ValueWitness(val)
	if err != nil {
		return errors.Wrapf(err, "getting composite %q value witness", s.Name())
	}
	exist, err := ss.internalHasValue(we.Key)
	if err != nil {
		return errors.Wrapf(err, "checking composite %q witness existence", s.Name())
	}
	if !exist {
		if err := ss.internalPutValue(we.Key, we.Value); err != nil {
			return errors.Wrapf(err, "putting composite %q witness", s.Name())
		}
	}
	hascomps := false
	entries, err := s.SingletonsEntries(val)
	if err != nil {
		return errors.WithStack(err)
	}
	if len(entries) > 0 {
		hascomps = true
	}
	for _, entry := range entries {
		if !reflect.ValueOf(entry.Value).IsNil() {
			if err := ss.internalPutValue(entry.Key, entry.Value); err != nil {
				return errors.Wrapf(err, "putting composite %q singleton %q", s.Name(), entry)
			}
		}
	}
	entries, err = s.CollectionsEntries(val)
	if err != nil {
		return errors.WithStack(err)
	}
	if len(entries) > 0 {
		hascomps = true
	}
	for _, entry := range entries {
		if reflect.ValueOf(entry.Value).IsNil() {
			if err := ss.internalDelValue(entry.Key); err != nil {
				return errors.Wrapf(err, "deleting composite %q collection entry %q", s.Name(), entry)
			}
		} else {
			if err := ss.internalPutValue(entry.Key, entry.Value); err != nil {
				return errors.Wrapf(err, "putting composite %q collection entry %q", s.Name(), entry)
			}
		}
	}
	if !hascomps || s.MustKeepRoot(val) {
		entry, err := s.RootEntry(val)
		if err != nil {
			return errors.WithStack(err)
		}
		if err := ss.internalPutValue(entry.Key, entry.Value); err != nil {
			return errors.Wrapf(err, "putting composite %q root entry %q", s.Name(), entry)
		}
	}
	return nil
}

func (ss *simplestore) GetComposite(s *Schema, id interface{}) (interface{}, error) {
	valkey, err := s.IdentifierKey(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if ok, err := ss.HasComposite(s, id); err != nil {
		return nil, errors.Wrapf(err, "checking composite %q with key %q existence", s.Name(), valkey)
	} else if !ok {
		return nil, nil // no existe la persona
	}
	val, err := s.Create()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = s.SetIdentifier(val, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	states, err := ss.stub.GetStateByRange(valkey.RangeUsing(ss.sep))
	if err != nil {
		return nil, errors.Wrapf(err, "getting composite %q with key %q states iterator", s.Name(), valkey)
	}
	merrs := []MemberError{}
	defer states.Close()
	for states.HasNext() {
		state, err := states.Next()
		if err != nil {
			return nil, errors.Wrapf(err, "getting composite %q with key %q next state", s.Name(), valkey)
		}
		statekey, err := key.ParseUsing(state.GetKey(), ss.sep)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing composite %q with key %q item", s.Name(), state.GetKey())
		}
		merr := ss.inject(s, statekey, state, valkey, val)
		if ss.seterrs && merr != nil {
			merrs = append(merrs, *merr)
		}
	}
	if ss.seterrs && len(merrs) > 0 {
		seterrs(val, merrs)
	}
	return val, nil
}

func (ss *simplestore) HasComposite(s *Schema, id interface{}) (bool, error) {
	key, err := s.IdentifierKey(id)
	if err != nil {
		return false, errors.WithStack(err)
	}
	wk := s.KeyWitness(key)
	var a interface{}
	found, err := ss.internalGetValue(wk, &a)
	if err != nil {
		return false, errors.Wrapf(err, "getting composite %q witness with key %q", s.Name(), wk.StringUsing(ss.sep))
	}
	return found, nil
}

func (ss *simplestore) DelComposite(s *Schema, id interface{}) error {
	key, err := s.IdentifierKey(id)
	if err != nil {
		return errors.WithStack(err)
	}
	states, err := ss.stub.GetStateByRange(key.RangeUsing(ss.sep))
	if err != nil {
		return errors.Wrapf(err, "getting composite %q states with key %q for deletion", s.Name(), key)
	}
	defer states.Close()
	for states.HasNext() {
		state, err := states.Next()
		if err != nil {
			return errors.Wrapf(err, "getting composite %q with key %q next state for deletion", s.Name(), key)
		}
		err = ss.stub.DelState(state.GetKey())
		if err != nil {
			return errors.Wrapf(err, "deleting composite %q with key %q state %q", s.Name(), key, state.GetKey())
		}
	}
	return nil
}

func (ss *simplestore) DelCompositeRange(s *Schema, r *Range) ([]interface{}, error) {
	first, last, err := ss.identifierKeyRange(s, r)
	if err != nil {
		return nil, errors.Wrapf(err, "getting keys range %v", r)
	}
	states, err := ss.stub.GetStateByRange(first, last)
	if err != nil {
		return nil, errors.Wrapf(err, "getting composite %q range [%q,%q] for deletion", s.Name(), first, last)
	}
	defer states.Close()
	res := []interface{}{}
	for states.HasNext() {
		state, err := states.Next()
		if err != nil {
			return nil, errors.Wrapf(err, "getting composite %q range [%q,%q] next key for deletion", s.Name(), first, last)
		}
		statekey, err := key.ParseUsing(state.GetKey(), ss.sep)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing state key %q as composite %q key", state.GetKey(), s.Name())
		}
		if s.IsWitnessKey(statekey) {
			id, err := s.KeyIdentifier(statekey)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			res = append(res, id)
		}
		err = ss.stub.DelState(state.GetKey())
		if err != nil {
			return nil, errors.Wrapf(err, "deleting composite %q range [%q,%q] state %q", s.Name(), first, last, state.GetKey())
		}
	}
	return res, nil
}

func (ss *simplestore) GetCompositeRange(s *Schema, r *Range) ([]interface{}, error) {
	first, last, err := ss.identifierKeyRange(s, r)
	if err != nil {
		return nil, errors.Wrapf(err, "getting keys range %v", r)
	}
	states, err := ss.stub.GetStateByRange(first, last)
	if err != nil {
		return nil, errors.Wrapf(err, "getting composite %q range [%q,%q] for reading", s.Name(), first, last)
	}
	defer states.Close()
	return ss.internalReadCompositeIterator(s, states)
}

func (ss *simplestore) GetCompositeAll(s *Schema) ([]interface{}, error) {
	kbn := s.KeyBaseName()
	if kbn == "" {
		return nil, errors.Errorf("getting composite %q all instances: keybasename is empty", s.Name())
	}
	first, last := key.NewBase(kbn, "").RangeUsing(ss.sep)
	states, err := ss.stub.GetStateByRange(first, last)
	if err != nil {
		return nil, errors.Wrapf(err, "getting composite %q all instances for reading", s.Name())
	}
	defer states.Close()
	return ss.internalReadCompositeIterator(s, states)
}

func (ss *simplestore) internalReadCompositeIterator(s *Schema, states shim.StateQueryIteratorInterface) ([]interface{}, error) {
	var (
		valkey  *key.Key
		val, id interface{}
	)
	merrs := []MemberError{}
	res := []interface{}{}
	for states.HasNext() {
		state, err := states.Next()
		if err != nil {
			return nil, errors.Wrapf(err, "getting composite %q iterator next key for reading", s.Name())
		}
		statekey, err := key.ParseUsing(state.GetKey(), ss.sep)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing state key %q as composite %q key", state.GetKey(), s.Name())
		}
		basekey := key.NewBaseKey(statekey)
		if valkey == nil || !valkey.Equal(basekey) {
			if ss.seterrs && len(merrs) > 0 {
				seterrs(val, merrs)
				merrs = []MemberError{}
			}
			valkey = basekey
			val, err = s.Create()
			if err != nil {
				return nil, errors.WithStack(err)
			}
			id, err = s.KeyIdentifier(valkey)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			err = s.SetIdentifier(val, id)
			if err != nil {
				return nil, errors.Wrapf(err, "setting composite %q id %v from key %v", s.Name(), id, valkey)
			}
			res = append(res, val)
		}
		merr := ss.inject(s, statekey, state, valkey, val)
		if ss.seterrs && merr != nil {
			merrs = append(merrs, *merr)
		}
	}
	return res, nil
}

func (ss *simplestore) internalPutValue(k *key.Key, value interface{}) error {
	if err := k.ValidateUsing(ss.sep); err != nil {
		return errors.Wrap(err, "checking value key")
	} else if bs, err := ss.marshaling.Marshal(value); err != nil {
		return errors.Wrap(err, "marshaling value")
	} else if bs, err := ss.filtering.Filter(bs); err != nil {
		return errors.Wrap(err, "filtering value")
	} else {
		ks := k.StringUsing(ss.sep)
		if log.IsEnabledFor(shim.LogDebug) {
			log.Debugf("putting key '%s' with value '%s'", ks, string(bs))
		}
		if err := ss.stub.PutState(ks, bs); err != nil {
			return errors.Wrap(err, "putting marshaled value into state")
		}
	}
	return nil
}

func (ss *simplestore) internalHasValue(k *key.Key) (bool, error) {
	bs, err := ss.stub.GetState(k.StringUsing(ss.sep))
	if err != nil {
		return false, errors.Wrap(err, "getting value from state")
	}
	return bs != nil, nil
}

func (ss *simplestore) internalGetValue(k *key.Key, value interface{}) (bool, error) {
	bs, err := ss.stub.GetState(k.StringUsing(ss.sep))
	if err != nil {
		return false, errors.Wrap(err, "getting marshaled value from state")
	}
	if bs == nil {
		return false, nil
	}
	return true, ss.internalParseValue(bs, value)
}

func (ss *simplestore) internalDelValue(k *key.Key) error {
	err := ss.stub.DelState(k.StringUsing(ss.sep))
	if err != nil {
		return errors.Wrap(err, "deleting value from state")
	}
	return nil
}

func (ss *simplestore) internalParseValue(bs []byte, value interface{}) error {
	if bs, err := ss.filtering.Unfilter(bs); err != nil {
		return errors.Wrap(err, "unfiltering value")
	} else if err := ss.marshaling.Unmarshal(bs, value); err != nil {
		return errors.Wrap(err, "unmarshaling value")
	}
	return nil
}

func (ss *simplestore) inject(s *Schema, statekey *key.Key, state *queryresult.KV, valkey *key.Key, val interface{}) *MemberError {
	var merr *MemberError
	switch {
	case statekey.Equal(valkey):
		err := ss.internalParseValue(state.GetValue(), val)
		if err != nil {
			ss.log.Errorf("parsing composite %q with key root item %q value in tx %s: %v", s.Name, valkey, ss.stub.GetTxID(), err)
			if ss.seterrs {
				seterr(val, err)
			}
			merr = &MemberError{
				Kind:  "root",
				Error: err.Error(),
			}
		}
	case s.Collection(statekey) != nil:
		member := s.Collection(statekey)
		itemval := member.Creator()
		err := ss.internalParseValue(state.GetValue(), itemval)
		if err != nil {
			ss.log.Errorf("parsing composite %q with key %q collection item %q value in tx %s: %v", s.Name, valkey, statekey, ss.stub.GetTxID(), err)
			if ss.seterrs {
				seterr(itemval, err)
			}
			merr = &MemberError{
				Kind:  "collection",
				Tag:   member.Tag,
				ID:    statekey.Tag.Value,
				Error: err.Error(),
			}
		}
		member.Collector(val, Item{Identifier: statekey.Tag.Value, Value: itemval})
	case s.Singleton(statekey) != nil:
		member := s.Singleton(statekey)
		itemval := member.Creator()
		err := ss.internalParseValue(state.GetValue(), itemval)
		if err != nil {
			ss.log.Errorf("parsing composite %q with key %q collection item %q value in tx %s: %v", s.Name, valkey, statekey, ss.stub.GetTxID(), err)
			if ss.seterrs {
				seterr(itemval, err)
			}
			merr = &MemberError{
				Kind:  "singleton",
				Tag:   member.Tag,
				Error: err.Error(),
			}
		}
		member.Setter(val, itemval)
	}
	return merr
}

func (ss *simplestore) identifierKeyRange(s *Schema, r *Range) (string, string, error) {
	fk, err := s.IdentifierKey(r.First)
	if err != nil {
		return "", "", errors.Wrap(err, "getting range start key")
	}
	lk, err := s.IdentifierKey(r.Last)
	if err != nil {
		return "", "", errors.Wrap(err, "getting range end key")
	}
	first, _ := fk.RangeUsing(ss.sep)
	_, last := lk.RangeUsing(ss.sep)
	return first, last, nil
}

func seterr(val interface{}, e error) {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		f := v.FieldByName("Error")
		if f.Kind() == reflect.String {
			f.SetString(e.Error())
		}
	}
}

func seterrs(val interface{}, merrs []MemberError) {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		f := v.FieldByName("Errors")
		if f.Kind() == reflect.Interface {
			f.Set(reflect.ValueOf(merrs))
		}
	}
}
