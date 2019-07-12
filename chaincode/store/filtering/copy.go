package filtering

func Copy() Filtering {
	return New(FilterFunc(copy), UnfilterFunc(copy))
}

func copy(bs []byte) ([]byte, error) {
	return bs, nil
}
