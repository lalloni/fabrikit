package key

import (
	"github.com/pkg/errors"
)

const (
	parsingBaseName int = iota
	parsingBaseValue
	parsingTagName
	parsingTagValue
)

func Parse(s string) (*Key, error) {
	return ParseUsing(s, DefaultSep)
}

func ParseUsing(s string, sep *Sep) (*Key, error) {
	var (
		base  = []Seg{{}}
		seg   = &base[0]
		tag   = Seg{}
		state = parsingBaseName
	)
	for i, r := range s {
		switch state {
		case parsingBaseName:
			if sep.IsAny(r) {
				if len(seg.Name) == 0 {
					return nil, errors.Errorf("illegal separator %q at position %d in %q: expecting segment name", r, i+1, s)
				}
				if r != sep.NameValueSep {
					return nil, errors.Errorf("illegal separator %q at position %d in %q: expecting name value separator %q", r, i+1, s, sep.NameValueSep)
				}
				state = parsingBaseValue
				continue
			}
			seg.Name += string(r)
		case parsingBaseValue:
			if sep.IsAny(r) {
				if len(seg.Value) == 0 {
					return nil, errors.Errorf("illegal separator %q at position %d in %q: expecting segment value", r, i+1, s)
				}
				switch r {
				case sep.TagSep:
					state = parsingTagName
				case sep.SegSep:
					base = append(base, Seg{})
					seg = &base[len(base)-1]
					state = parsingBaseName
				case sep.NameValueSep:
					return nil, errors.Errorf("illegal separator %q at position %d in %q: expecting segment separator %q or tag separator %q", r, i+1, s, sep.SegSep, sep.TagSep)
				}
				continue
			}
			seg.Value += string(r)
		case parsingTagName:
			if sep.IsAny(r) {
				if len(tag.Name) == 0 {
					return nil, errors.Errorf("illegal separator %q at position %d in %q: expecting segment name", r, i+1, s)
				}
				if r != sep.NameValueSep {
					return nil, errors.Errorf("illegal separator %q at position %d in %q: expecting name value separator %q", r, i+1, s, sep.NameValueSep)
				}
				state = parsingTagValue
				continue
			}
			tag.Name += string(r)
		case parsingTagValue:
			if sep.IsAny(r) {
				return nil, errors.Errorf("illegal separator %q at position %d in %q: expecting tag value", r, i+1, s)
			}
			tag.Value += string(r)
		}
	}
	if state == parsingTagName && len(tag.Name) == 0 {
		return nil, errors.Errorf("illegal end at position %d in %q: expecting tag name", len(s), s)
	}
	return &Key{
		Base: base,
		Tag:  tag,
	}, nil
}
