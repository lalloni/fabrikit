package filtering

import (
	"bytes"
	"io"

	"github.com/pierrec/lz4"
)

func LZ4() Filtering {
	return New(
		FilterFunc(func(bs []byte) ([]byte, error) {
			buf := &bytes.Buffer{}
			w := lz4.NewWriter(buf)
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
			r := lz4.NewReader(bytes.NewReader(bs))
			if _, err := io.Copy(buf, r); err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}))
}
