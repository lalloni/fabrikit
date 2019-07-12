package key_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lalloni/fabrikit/chaincode/store/key"
)

func TestParseUsing(t *testing.T) {
	sep := &key.Sep{
		NameValueSep: '=',
		SegSep:       ';',
		TagSep:       '#',
	}
	tests := []struct {
		name string
		s    string
		err  bool
		k    *key.Key
	}{
		{"one-segment", "a=1", false, key.NewBase("a", "1")},
		{"two-segments-no-tag", "a=1;b=2", false, key.NewBase("a", "1", "b", "2")},
		{"two-segments-and-tag", "a=1;b=2#pepe", false, key.NewBase("a", "1", "b", "2").Tagged("pepe")},
		{"two-segments-and-tag-value", "a=1;b=2#pepe=3", false, key.NewBase("a", "1", "b", "2").Tagged("pepe", "3")},
		{"no-segment-value-1", "a#2", true, nil},
		{"no-segment-value-2", "a=#2", true, nil},
		{"no-segments", "#2", true, nil},
		{"empty-segment", ";a=1", true, nil},
		{"no-segment-name", "=2;a=1", true, nil},
		{"no-tagname", "a=1#", true, nil},
		{"no-segment-before-tag", "a=1;#", true, nil},
		{"two-segment-values", "a=1=2;#", true, nil},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			a := assert.New(t)
			k, err := key.ParseUsing(test.s, sep)
			if test.err {
				a.Error(err)
			} else {
				a.Equal(test.k, k)
			}
			t.Logf("test:%+v k:%+v err:%+v", test, k, err)
		})
	}
}
