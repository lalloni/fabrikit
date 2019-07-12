package key

import (
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
)

var DefaultSep = &Sep{
	NameValueSep: ':',
	SegSep:       '/',
	TagSep:       '#',
}

type Sep struct {
	// NameValueSep es el separador por defecto de nombres y valores claves
	NameValueSep rune
	// SegSep es el separador por defecto de segmentos de claves
	SegSep rune
	// TagSep es el separador por defecto de etiqueta de clave
	TagSep rune
}

func (s *Sep) IsAny(r rune) bool {
	return r == s.NameValueSep || r == s.SegSep || r == s.TagSep
}

// Key es una state key inmutable
type Key struct {
	Base []Seg
	Tag  Seg
}

func NewBase(ss ...string) *Key {
	s := len(ss)
	l := s - s%2
	k := &Key{}
	for i := 0; i < l; i += 2 {
		k.Base = append(k.Base, Seg{
			Name:  ss[i],
			Value: ss[i+1],
		})
	}
	return k
}

func NewBaseKey(k *Key) *Key {
	return &Key{Base: k.Base}
}

func (k *Key) Tagged(name string, value ...string) *Key {
	k2 := &Key{
		Base: k.Base,
		Tag:  Seg{Name: name},
	}
	if len(value) > 0 {
		k2.Tag.Value = value[0]
	}
	return k2
}

func (k *Key) AppendBase(name, value string) *Key {
	return &Key{
		Base: append(k.Base, Seg{Name: name, Value: value}),
		Tag:  k.Tag,
	}
}

func (k *Key) Equal(o *Key) bool {
	if len(k.Base) != len(o.Base) || k.Tag != o.Tag {
		return false
	}
	for i := 0; i < len(k.Base); i++ {
		if k.Base[i] != o.Base[i] {
			return false
		}
	}
	return true
}

func (k *Key) Validate() error {
	return k.ValidateUsing(DefaultSep)
}

func (k *Key) ValidateUsing(sep *Sep) error {
	for _, seg := range k.Base {
		seg := seg
		if err := validateSeg(&seg, sep); err != nil {
			return err
		}
	}
	if err := validateSeg(&k.Tag, sep); err != nil {
		return err
	}
	return nil
}

func (k *Key) String() string {
	return k.StringUsing(DefaultSep)
}

func (k *Key) StringUsing(sep *Sep) string {
	s := &strings.Builder{}
	l := len(k.Base) - 1
	for i, seg := range k.Base {
		s.WriteString(seg.StringUsing(sep))
		if i < l {
			s.WriteRune(sep.SegSep)
		}
	}
	if len(k.Tag.Name) > 0 {
		s.WriteRune(sep.TagSep)
		s.WriteString(k.Tag.StringUsing(sep))
	}
	return s.String()
}

func (k *Key) Range() (string, string) {
	s := k.String()
	return s, s + string(utf8.MaxRune)
}

func (k *Key) RangeUsing(sep *Sep) (string, string) {
	s := k.StringUsing(sep)
	return s, s + string(utf8.MaxRune)
}

// Seg es par clave-valor de una key
type Seg struct {
	Name  string
	Value string
}

func (s *Seg) String() string {
	return s.StringUsing(DefaultSep)
}

func (s *Seg) StringUsing(sep *Sep) string {
	if len(s.Name) == 0 {
		return ""
	}
	b := &strings.Builder{}
	b.WriteString(s.Name)
	if len(s.Value) > 0 {
		b.WriteRune(sep.NameValueSep)
		b.WriteString(s.Value)
	}
	return b.String()
}

func validateSeg(seg *Seg, sep *Sep) error {
	if err := validateString(&seg.Name, sep); err != nil {
		return errors.Wrapf(err, "checking key base segment %q name", seg)
	}
	if err := validateString(&seg.Value, sep); err != nil {
		return errors.Wrapf(err, "checking key base segment %q value", seg)
	}
	return nil
}

func validateString(str *string, sep *Sep) error {
	for _, r := range *str {
		if sep.IsAny(r) {
			return errors.Errorf("character not allowed %q", r)
		}
	}
	return nil
}
