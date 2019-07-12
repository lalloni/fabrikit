package store

import (
	"github.com/lalloni/fabrikit/chaincode/store/filtering"
	"github.com/lalloni/fabrikit/chaincode/store/key"
	"github.com/lalloni/fabrikit/chaincode/store/marshaling"
)

type Option func(*simplestore)

func SetSep(sep *key.Sep) Option {
	return func(s *simplestore) {
		s.sep = sep
	}
}

func SetMarshaling(m marshaling.Marshaling) Option {
	return func(s *simplestore) {
		s.marshaling = m
	}
}

func SetFiltering(f filtering.Filtering) Option {
	return func(s *simplestore) {
		s.filtering = f
	}
}

func SetErrors(b bool) Option {
	return func(s *simplestore) {
		s.seterrs = b
	}
}
