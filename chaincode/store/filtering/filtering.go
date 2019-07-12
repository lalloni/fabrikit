package filtering

type Filtering interface {
	Filter
	Unfilter
}

type Filter interface {
	Filter([]byte) ([]byte, error)
}

type Unfilter interface {
	Unfilter(bs []byte) ([]byte, error)
}

func New(f Filter, u Unfilter) Filtering {
	return filtering{f, u}
}

type filtering struct {
	f Filter
	u Unfilter
}

func (m filtering) Filter(bs []byte) ([]byte, error) {
	return m.f.Filter(bs)
}

func (m filtering) Unfilter(bs []byte) ([]byte, error) {
	return m.u.Unfilter(bs)
}

type FilterFunc func([]byte) ([]byte, error)

func (m FilterFunc) Filter(bs []byte) ([]byte, error) {
	return m(bs)
}

type UnfilterFunc func([]byte) ([]byte, error)

func (f UnfilterFunc) Unfilter(bs []byte) ([]byte, error) {
	return f(bs)
}
