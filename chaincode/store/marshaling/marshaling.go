package marshaling

type Marshaling interface {
	Marshaler
	Unmarshaler
}

type Marshaler interface {
	Marshal(value interface{}) ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal(bs []byte, value interface{}) error
}

func New(m Marshaler, u Unmarshaler) Marshaling {
	return marshaling{m, u}
}

type marshaling struct {
	m Marshaler
	u Unmarshaler
}

func (m marshaling) Marshal(value interface{}) ([]byte, error) {
	return m.m.Marshal(value)
}

func (m marshaling) Unmarshal(bs []byte, value interface{}) error {
	return m.u.Unmarshal(bs, value)
}

type MarshalerFunc func(interface{}) ([]byte, error)

func (m MarshalerFunc) Marshal(v interface{}) ([]byte, error) {
	return m(v)
}

type UnmarshalerFunc func([]byte, interface{}) error

func (f UnmarshalerFunc) Unmarshal(bs []byte, value interface{}) error {
	return f(bs, value)
}
