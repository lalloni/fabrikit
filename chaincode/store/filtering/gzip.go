package filtering

import (
	"bytes"
	"compress/gzip"
	"io"
)

func GZip() Filtering {
	return New(
		FilterFunc(func(bs []byte) ([]byte, error) {
			buf := &bytes.Buffer{}
			w := gzip.NewWriter(buf)
			if _, err := io.Copy(w, bytes.NewReader(bs)); err != nil {
				return nil, err
			}
			if err := w.Close(); err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}),
		UnfilterFunc(func(bs []byte) ([]byte, error) {
			buf := &bytes.Buffer{}
			r, err := gzip.NewReader(bytes.NewReader(bs))
			if err != nil {
				return nil, err
			}
			defer r.Close()
			if _, err := io.Copy(buf, r); err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}))
}
